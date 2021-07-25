package alerting

import (
	"context"
	"fmt"

	"github.com/openinsight-project/grafinsight/pkg/bus"
	"github.com/openinsight-project/grafinsight/pkg/components/simplejson"
	"github.com/openinsight-project/grafinsight/pkg/models"
)

// AlertTestCommand initiates an test evaluation
// of an alert rule.
type AlertTestCommand struct {
	Dashboard *simplejson.Json
	PanelID   int64
	OrgID     int64
	User      *models.SignedInUser

	Result *EvalContext
}

func init() {
	bus.AddHandler("alerting", handleAlertTestCommand)
}

func handleAlertTestCommand(cmd *AlertTestCommand) error {
	dash := models.NewDashboardFromJson(cmd.Dashboard)

	extractor := NewDashAlertExtractor(dash, cmd.OrgID, cmd.User)
	alerts, err := extractor.GetAlerts()
	if err != nil {
		return err
	}

	for _, alert := range alerts {
		if alert.PanelId == cmd.PanelID {
			rule, err := NewRuleFromDBAlert(alert, true)
			if err != nil {
				return err
			}

			cmd.Result = testAlertRule(rule)
			return nil
		}
	}

	return fmt.Errorf("could not find alert with panel ID %d", cmd.PanelID)
}

func testAlertRule(rule *Rule) *EvalContext {
	handler := NewEvalHandler()

	context := NewEvalContext(context.Background(), rule, fakeRequestValidator{})
	context.IsTestRun = true
	context.IsDebug = true

	handler.Eval(context)
	context.Rule.State = context.GetNewState()

	return context
}
