package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/openinsight-project/grafinsight/pkg/infra/log"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/setting"
	macaron "gopkg.in/macaron.v1"
)

// AddCSPHeader adds the Content Security Policy header.
func AddCSPHeader(cfg *setting.Cfg, logger log.Logger) macaron.Handler {
	return func(w http.ResponseWriter, req *http.Request, c *macaron.Context) {
		if !cfg.CSPEnabled {
			return
		}

		logger.Debug("Adding CSP header to response", "cfg", fmt.Sprintf("%p", cfg))

		ctx, ok := c.Data["ctx"].(*models.ReqContext)
		if !ok {
			panic("Failed to convert context into models.ReqContext")
		}

		if cfg.CSPTemplate == "" {
			logger.Debug("CSP template not configured, so returning 500")
			ctx.JsonApiErr(500, "CSP template has to be configured", nil)
			return
		}

		var buf [16]byte
		if _, err := io.ReadFull(rand.Reader, buf[:]); err != nil {
			logger.Error("Failed to generate CSP nonce", "err", err)
			ctx.JsonApiErr(500, "Failed to generate CSP nonce", err)
		}

		nonce := base64.RawStdEncoding.EncodeToString(buf[:])
		val := strings.ReplaceAll(cfg.CSPTemplate, "$NONCE", fmt.Sprintf("'nonce-%s'", nonce))
		w.Header().Set("Content-Security-Policy", val)
		ctx.RequestNonce = nonce
		logger.Debug("Successfully generated CSP nonce", "nonce", nonce)
	}
}
