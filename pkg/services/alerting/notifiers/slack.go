package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/openinsight-project/grafinsight/pkg/bus"
	"github.com/openinsight-project/grafinsight/pkg/infra/log"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/services/alerting"
	"github.com/openinsight-project/grafinsight/pkg/setting"
)

func init() {
	alerting.RegisterNotifier(&alerting.NotifierPlugin{
		Type:        "slack",
		Name:        "Slack",
		Description: "Sends notifications to Slack via Slack Webhooks",
		Heading:     "Slack settings",
		Factory:     NewSlackNotifier,
		Options: []alerting.NotifierOption{
			{
				Label:        "Url",
				Element:      alerting.ElementTypeInput,
				InputType:    alerting.InputTypeText,
				Placeholder:  "Slack incoming webhook url",
				PropertyName: "url",
				Required:     true,
				Secure:       true,
			},
			{
				Label:        "Recipient",
				Element:      alerting.ElementTypeInput,
				InputType:    alerting.InputTypeText,
				Description:  "Override default channel or user, use #channel-name, @username (has to be all lowercase, no whitespace), or user/channel Slack ID",
				PropertyName: "recipient",
			},
			{
				Label:        "Username",
				Element:      alerting.ElementTypeInput,
				InputType:    alerting.InputTypeText,
				Description:  "Set the username for the bot's message",
				PropertyName: "username",
			},
			{
				Label:        "Icon emoji",
				Element:      alerting.ElementTypeInput,
				InputType:    alerting.InputTypeText,
				Description:  "Provide an emoji to use as the icon for the bot's message. Overrides the icon URL.",
				PropertyName: "iconEmoji",
			},
			{
				Label:        "Icon URL",
				Element:      alerting.ElementTypeInput,
				InputType:    alerting.InputTypeText,
				Description:  "Provide a URL to an image to use as the icon for the bot's message",
				PropertyName: "iconUrl",
			},
			{
				Label:        "Mention Users",
				Element:      alerting.ElementTypeInput,
				InputType:    alerting.InputTypeText,
				Description:  "Mention one or more users (comma separated) when notifying in a channel, by ID (you can copy this from the user's Slack profile)",
				PropertyName: "mentionUsers",
			},
			{
				Label:        "Mention Groups",
				Element:      alerting.ElementTypeInput,
				InputType:    alerting.InputTypeText,
				Description:  "Mention one or more groups (comma separated) when notifying in a channel (you can copy this from the group's Slack profile URL)",
				PropertyName: "mentionGroups",
			},
			{
				Label:   "Mention Channel",
				Element: alerting.ElementTypeSelect,
				SelectOptions: []alerting.SelectOption{
					{
						Value: "",
						Label: "Disabled",
					},
					{
						Value: "here",
						Label: "Every active channel member",
					},
					{
						Value: "channel",
						Label: "Every channel member",
					},
				},
				Description:  "Mention whole channel or just active members when notifying",
				PropertyName: "mentionChannel",
			},
			{
				Label:        "Token",
				Element:      alerting.ElementTypeInput,
				InputType:    alerting.InputTypeText,
				Description:  "Provide a bot token to use the Slack file.upload API (starts with \"xoxb\"). Specify Recipient for this to work",
				PropertyName: "token",
				Secure:       true,
			},
		},
	})
}

var reRecipient *regexp.Regexp = regexp.MustCompile("^((@[a-z0-9][a-zA-Z0-9._-]*)|(#[^ .A-Z]{1,79})|([a-zA-Z0-9]+))$")

// NewSlackNotifier is the constructor for the Slack notifier
func NewSlackNotifier(model *models.AlertNotification) (alerting.Notifier, error) {
	url := model.DecryptedValue("url", model.Settings.Get("url").MustString())
	if url == "" {
		return nil, alerting.ValidationError{Reason: "Could not find url property in settings"}
	}

	recipient := strings.TrimSpace(model.Settings.Get("recipient").MustString())
	if recipient != "" && !reRecipient.MatchString(recipient) {
		return nil, alerting.ValidationError{Reason: fmt.Sprintf("Recipient on invalid format: %q", recipient)}
	}
	username := model.Settings.Get("username").MustString()
	iconEmoji := model.Settings.Get("icon_emoji").MustString()
	iconURL := model.Settings.Get("icon_url").MustString()
	mentionUsersStr := model.Settings.Get("mentionUsers").MustString()
	mentionGroupsStr := model.Settings.Get("mentionGroups").MustString()
	mentionChannel := model.Settings.Get("mentionChannel").MustString()
	token := model.DecryptedValue("token", model.Settings.Get("token").MustString())

	uploadImage := model.Settings.Get("uploadImage").MustBool(true)

	if mentionChannel != "" && mentionChannel != "here" && mentionChannel != "channel" {
		return nil, alerting.ValidationError{
			Reason: fmt.Sprintf("Invalid value for mentionChannel: %q", mentionChannel),
		}
	}
	mentionUsers := []string{}
	for _, u := range strings.Split(mentionUsersStr, ",") {
		u = strings.TrimSpace(u)
		if u != "" {
			mentionUsers = append(mentionUsers, u)
		}
	}
	mentionGroups := []string{}
	for _, g := range strings.Split(mentionGroupsStr, ",") {
		g = strings.TrimSpace(g)
		if g != "" {
			mentionGroups = append(mentionGroups, g)
		}
	}

	return &SlackNotifier{
		NotifierBase:   NewNotifierBase(model),
		URL:            url,
		Recipient:      recipient,
		Username:       username,
		IconEmoji:      iconEmoji,
		IconURL:        iconURL,
		MentionUsers:   mentionUsers,
		MentionGroups:  mentionGroups,
		MentionChannel: mentionChannel,
		Token:          token,
		Upload:         uploadImage,
		log:            log.New("alerting.notifier.slack"),
	}, nil
}

// SlackNotifier is responsible for sending
// alert notification to Slack.
type SlackNotifier struct {
	NotifierBase
	URL            string
	Recipient      string
	Username       string
	IconEmoji      string
	IconURL        string
	MentionUsers   []string
	MentionGroups  []string
	MentionChannel string
	Token          string
	Upload         bool
	log            log.Logger
}

// Notify send alert notification to Slack.
func (sn *SlackNotifier) Notify(evalContext *alerting.EvalContext) error {
	sn.log.Info("Executing slack notification", "ruleId", evalContext.Rule.ID, "notification", sn.Name)

	ruleURL, err := evalContext.GetRuleURL()
	if err != nil {
		sn.log.Error("Failed get rule link", "error", err)
		return err
	}

	fields := make([]map[string]interface{}, 0)
	fieldLimitCount := 4
	for index, evt := range evalContext.EvalMatches {
		fields = append(fields, map[string]interface{}{
			"title": evt.Metric,
			"value": evt.Value,
			"short": true,
		})
		if index > fieldLimitCount {
			break
		}
	}

	if evalContext.Error != nil {
		fields = append(fields, map[string]interface{}{
			"title": "Error message",
			"value": evalContext.Error.Error(),
			"short": false,
		})
	}

	mentionsBuilder := strings.Builder{}
	appendSpace := func() {
		if mentionsBuilder.Len() > 0 {
			mentionsBuilder.WriteString(" ")
		}
	}
	mentionChannel := strings.TrimSpace(sn.MentionChannel)
	if mentionChannel != "" {
		mentionsBuilder.WriteString(fmt.Sprintf("<!%s|%s>", mentionChannel, mentionChannel))
	}
	if len(sn.MentionGroups) > 0 {
		appendSpace()
		for _, g := range sn.MentionGroups {
			mentionsBuilder.WriteString(fmt.Sprintf("<!subteam^%s>", g))
		}
	}
	if len(sn.MentionUsers) > 0 {
		appendSpace()
		for _, u := range sn.MentionUsers {
			mentionsBuilder.WriteString(fmt.Sprintf("<@%s>", u))
		}
	}
	msg := ""
	if evalContext.Rule.State != models.AlertStateOK { // don't add message when going back to alert state ok.
		msg = evalContext.Rule.Message
	}
	imageURL := ""
	// default to file.upload API method if a token is provided
	if sn.Token == "" {
		imageURL = evalContext.ImagePublicURL
	}

	var blocks []map[string]interface{}
	if mentionsBuilder.Len() > 0 {
		blocks = []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": mentionsBuilder.String(),
				},
			},
		}
	}
	attachment := map[string]interface{}{
		"color":       evalContext.GetStateModel().Color,
		"title":       evalContext.GetNotificationTitle(),
		"title_link":  ruleURL,
		"text":        msg,
		"fallback":    evalContext.GetNotificationTitle(),
		"fields":      fields,
		"footer":      "Grafinsight v" + setting.BuildVersion,
		"footer_icon": "https://grafinsight.com/assets/img/fav32.png",
		"ts":          time.Now().Unix(),
	}
	if sn.NeedsImage() && imageURL != "" {
		attachment["image_url"] = imageURL
	}
	body := map[string]interface{}{
		"text": evalContext.GetNotificationTitle(),
		"attachments": []map[string]interface{}{
			attachment,
		},
		"parse": "full", // to linkify urls, users and channels in alert message.
	}
	if len(blocks) > 0 {
		body["blocks"] = blocks
	}

	// recipient override
	if sn.Recipient != "" {
		body["channel"] = sn.Recipient
	}
	if sn.Username != "" {
		body["username"] = sn.Username
	}
	if sn.IconEmoji != "" {
		body["icon_emoji"] = sn.IconEmoji
	}
	if sn.IconURL != "" {
		body["icon_url"] = sn.IconURL
	}
	data, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	cmd := &models.SendWebhookSync{
		Url:        sn.URL,
		Body:       string(data),
		HttpMethod: http.MethodPost,
	}
	if sn.Token != "" {
		sn.log.Debug("Adding authorization header to HTTP request")
		cmd.HttpHeader = map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", sn.Token),
		}
	}
	if err := bus.DispatchCtx(evalContext.Ctx, cmd); err != nil {
		sn.log.Error("Failed to send slack notification", "error", err, "webhook", sn.Name)
		return err
	}
	if sn.Token != "" && sn.UploadImage {
		err = sn.slackFileUpload(evalContext, sn.log, "https://slack.com/api/files.upload", sn.Recipient, sn.Token)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sn *SlackNotifier) slackFileUpload(evalContext *alerting.EvalContext, log log.Logger, url string, recipient string, token string) error {
	if evalContext.ImageOnDiskPath == "" {
		// nolint:gosec
		// We can ignore the gosec G304 warning on this one because `setting.HomePath` comes from Grafinsight's configuration file.
		evalContext.ImageOnDiskPath = filepath.Join(setting.HomePath, "public/img/mixed_styles.png")
	}
	log.Info("Uploading to slack via file.upload API")
	headers, uploadBody, err := sn.generateSlackBody(evalContext.ImageOnDiskPath, token, recipient)
	if err != nil {
		return err
	}
	cmd := &models.SendWebhookSync{Url: url, Body: uploadBody.String(), HttpHeader: headers, HttpMethod: "POST"}
	if err := bus.DispatchCtx(evalContext.Ctx, cmd); err != nil {
		log.Error("Failed to upload slack image", "error", err, "webhook", "file.upload")
		return err
	}
	return nil
}

func (sn *SlackNotifier) generateSlackBody(path string, token string, recipient string) (map[string]string, bytes.Buffer, error) {
	// Slack requires all POSTs to files.upload to present
	// an "application/x-www-form-urlencoded" encoded querystring
	// See https://api.slack.com/methods/files.upload
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	defer func() {
		if err := w.Close(); err != nil {
			// Shouldn't matter since we already close w explicitly on the non-error path
			sn.log.Warn("Failed to close multipart writer", "err", err)
		}
	}()

	// Add the generated image file
	// We can ignore the gosec G304 warning on this one because `imagePath` comes
	// from the alert `evalContext` that generates the images. `evalContext` in turn derives the root of the file
	// path from configuration variables.
	// nolint:gosec
	f, err := os.Open(path)
	if err != nil {
		return nil, b, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			sn.log.Warn("Failed to close file", "path", path, "err", err)
		}
	}()
	fw, err := w.CreateFormFile("file", path)
	if err != nil {
		return nil, b, err
	}
	if _, err := io.Copy(fw, f); err != nil {
		return nil, b, err
	}
	// Add the authorization token
	if err := w.WriteField("token", token); err != nil {
		return nil, b, err
	}
	// Add the channel(s) to POST to
	if err := w.WriteField("channels", recipient); err != nil {
		return nil, b, err
	}
	if err := w.Close(); err != nil {
		return nil, b, fmt.Errorf("failed to close multipart writer: %w", err)
	}
	headers := map[string]string{
		"Content-Type":  w.FormDataContentType(),
		"Authorization": "auth_token=\"" + token + "\"",
	}
	return headers, b, nil
}
