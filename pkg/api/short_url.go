package api

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/openinsight-project/grafinsight/pkg/api/dtos"
	"github.com/openinsight-project/grafinsight/pkg/api/response"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/setting"
	"github.com/openinsight-project/grafinsight/pkg/util"
)

// createShortURL handles requests to create short URLs.
func (hs *HTTPServer) createShortURL(c *models.ReqContext, cmd dtos.CreateShortURLCmd) response.Response {
	hs.log.Debug("Received request to create short URL", "path", cmd.Path)

	cmd.Path = strings.TrimSpace(cmd.Path)

	if path.IsAbs(cmd.Path) {
		hs.log.Error("Invalid short URL path", "path", cmd.Path)
		return response.Error(400, "Path should be relative", nil)
	}

	shortURL, err := hs.ShortURLService.CreateShortURL(c.Req.Context(), c.SignedInUser, cmd.Path)
	if err != nil {
		return response.Error(500, "Failed to create short URL", err)
	}

	url := fmt.Sprintf("%s/goto/%s", strings.TrimSuffix(setting.AppUrl, "/"), shortURL.Uid)
	c.Logger.Debug("Created short URL", "url", url)

	dto := dtos.ShortURL{
		UID: shortURL.Uid,
		URL: url,
	}

	return response.JSON(200, dto)
}

func (hs *HTTPServer) redirectFromShortURL(c *models.ReqContext) {
	shortURLUID := c.Params(":uid")

	if !util.IsValidShortUID(shortURLUID) {
		return
	}

	shortURL, err := hs.ShortURLService.GetShortURLByUID(c.Req.Context(), c.SignedInUser, shortURLUID)
	if err != nil {
		if errors.Is(err, models.ErrShortURLNotFound) {
			hs.log.Debug("Not redirecting short URL since not found")
			return
		}

		hs.log.Error("Short URL redirection error", "err", err)
		return
	}

	// Failure to update LastSeenAt should still allow to redirect
	if err := hs.ShortURLService.UpdateLastSeenAt(c.Req.Context(), shortURL); err != nil {
		hs.log.Error("Failed to update short URL last seen at", "error", err)
	}

	hs.log.Debug("Redirecting short URL", "path", shortURL.Path)
	c.Redirect(setting.ToAbsUrl(shortURL.Path), 302)
}
