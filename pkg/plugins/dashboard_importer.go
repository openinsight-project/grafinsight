package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/openinsight-project/grafinsight/pkg/bus"
	"github.com/openinsight-project/grafinsight/pkg/components/simplejson"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/services/dashboards"
)

var varRegex = regexp.MustCompile(`(\$\{.+?\})`)

type ImportDashboardCommand struct {
	Dashboard *simplejson.Json
	Path      string
	Inputs    []ImportDashboardInput
	Overwrite bool
	FolderId  int64

	OrgId    int64
	User     *models.SignedInUser
	PluginId string
	Result   *PluginDashboardInfoDTO
}

type ImportDashboardInput struct {
	Type     string `json:"type"`
	PluginId string `json:"pluginId"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}

type DashboardInputMissingError struct {
	VariableName string
}

func (e DashboardInputMissingError) Error() string {
	return fmt.Sprintf("Dashboard input variable: %v missing from import command", e.VariableName)
}

func init() {
	bus.AddHandler("plugins", ImportDashboard)
}

func ImportDashboard(cmd *ImportDashboardCommand) error {
	var dashboard *models.Dashboard
	var err error

	if cmd.PluginId != "" {
		if dashboard, err = loadPluginDashboard(cmd.PluginId, cmd.Path); err != nil {
			return err
		}
	} else {
		dashboard = models.NewDashboardFromJson(cmd.Dashboard)
	}

	evaluator := &DashTemplateEvaluator{
		template: dashboard.Data,
		inputs:   cmd.Inputs,
	}

	generatedDash, err := evaluator.Eval()
	if err != nil {
		return err
	}

	saveCmd := models.SaveDashboardCommand{
		Dashboard: generatedDash,
		OrgId:     cmd.OrgId,
		UserId:    cmd.User.UserId,
		Overwrite: cmd.Overwrite,
		PluginId:  cmd.PluginId,
		FolderId:  cmd.FolderId,
	}

	dto := &dashboards.SaveDashboardDTO{
		OrgId:     cmd.OrgId,
		Dashboard: saveCmd.GetDashboardModel(),
		Overwrite: saveCmd.Overwrite,
		User:      cmd.User,
	}

	savedDash, err := dashboards.NewService().ImportDashboard(dto)

	if err != nil {
		return err
	}

	cmd.Result = &PluginDashboardInfoDTO{
		PluginId:         cmd.PluginId,
		Title:            savedDash.Title,
		Path:             cmd.Path,
		Revision:         savedDash.Data.Get("revision").MustInt64(1),
		FolderId:         savedDash.FolderId,
		ImportedUri:      "db/" + savedDash.Slug,
		ImportedUrl:      savedDash.GetUrl(),
		ImportedRevision: dashboard.Data.Get("revision").MustInt64(1),
		Imported:         true,
		DashboardId:      savedDash.Id,
		Slug:             savedDash.Slug,
	}

	return nil
}

type DashTemplateEvaluator struct {
	template  *simplejson.Json
	inputs    []ImportDashboardInput
	variables map[string]string
	result    *simplejson.Json
}

func (e *DashTemplateEvaluator) findInput(varName string, varType string) *ImportDashboardInput {
	for _, input := range e.inputs {
		if varType == input.Type && (input.Name == varName || input.Name == "*") {
			return &input
		}
	}

	return nil
}

func (e *DashTemplateEvaluator) Eval() (*simplejson.Json, error) {
	e.result = simplejson.New()
	e.variables = make(map[string]string)

	// check that we have all inputs we need
	for _, inputDef := range e.template.Get("__inputs").MustArray() {
		inputDefJson := simplejson.NewFromAny(inputDef)
		inputName := inputDefJson.Get("name").MustString()
		inputType := inputDefJson.Get("type").MustString()
		input := e.findInput(inputName, inputType)

		if input == nil {
			return nil, &DashboardInputMissingError{VariableName: inputName}
		}

		e.variables["${"+inputName+"}"] = input.Value
	}

	return simplejson.NewFromAny(e.evalObject(e.template)), nil
}

func (e *DashTemplateEvaluator) evalValue(source *simplejson.Json) interface{} {
	sourceValue := source.Interface()

	switch v := sourceValue.(type) {
	case string:
		interpolated := varRegex.ReplaceAllStringFunc(v, func(match string) string {
			replacement, exists := e.variables[match]
			if exists {
				return replacement
			}
			return match
		})
		return interpolated
	case bool:
		return v
	case json.Number:
		return v
	case map[string]interface{}:
		return e.evalObject(source)
	case []interface{}:
		array := make([]interface{}, 0)
		for _, item := range v {
			array = append(array, e.evalValue(simplejson.NewFromAny(item)))
		}
		return array
	}

	return nil
}

func (e *DashTemplateEvaluator) evalObject(source *simplejson.Json) interface{} {
	result := make(map[string]interface{})

	for key, value := range source.MustMap() {
		if key == "__inputs" {
			continue
		}
		result[key] = e.evalValue(simplejson.NewFromAny(value))
	}

	return result
}
