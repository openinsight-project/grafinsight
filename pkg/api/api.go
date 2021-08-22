// Package api contains API logic.
package api

import (
	"time"

	"github.com/go-macaron/binding"
	"github.com/openinsight-project/grafinsight/pkg/api/avatar"
	"github.com/openinsight-project/grafinsight/pkg/api/dtos"
	"github.com/openinsight-project/grafinsight/pkg/api/frontendlogging"
	"github.com/openinsight-project/grafinsight/pkg/api/routing"
	"github.com/openinsight-project/grafinsight/pkg/infra/log"
	"github.com/openinsight-project/grafinsight/pkg/middleware"
	"github.com/openinsight-project/grafinsight/pkg/models"
)

var plog = log.New("api")

// registerRoutes registers all API HTTP routes.
func (hs *HTTPServer) registerRoutes() {
	reqNoAuth := middleware.NoAuth()
	reqSignedIn := middleware.ReqSignedIn
	reqSignedInNoAnonymous := middleware.ReqSignedInNoAnonymous
	reqGrafinsightAdmin := middleware.ReqGrafinsightAdmin
	reqEditorRole := middleware.ReqEditorRole
	reqOrgAdmin := middleware.ReqOrgAdmin
	reqCanAccessTeams := middleware.AdminOrFeatureEnabled(hs.Cfg.EditorsCanAdmin)
	reqSnapshotPublicModeOrSignedIn := middleware.SnapshotPublicModeOrSignedIn(hs.Cfg)
	redirectFromLegacyDashboardURL := middleware.RedirectFromLegacyDashboardURL()
	redirectFromLegacyDashboardSoloURL := middleware.RedirectFromLegacyDashboardSoloURL(hs.Cfg)
	redirectFromLegacyPanelEditURL := middleware.RedirectFromLegacyPanelEditURL(hs.Cfg)
	quota := middleware.Quota(hs.QuotaService)
	bind := binding.Bind

	r := hs.RouteRegister

	// not logged in views
	r.Get("/logout", hs.Logout)
	r.Post("/login", quota("session"), bind(dtos.LoginCommand{}), routing.Wrap(hs.LoginPost))
	r.Get("/login/:name", quota("session"), hs.OAuthLogin)
	r.Get("/login", hs.LoginView)
	r.Get("/invite/:code", hs.Index)

	// authed views
	r.Get("/", reqSignedIn, hs.Index)
	r.Get("/profile/", reqSignedInNoAnonymous, hs.Index)
	r.Get("/profile/password", reqSignedInNoAnonymous, hs.Index)
	r.Get("/.well-known/change-password", redirectToChangePassword)
	r.Get("/profile/switch-org/:id", reqSignedInNoAnonymous, hs.ChangeActiveOrgAndRedirectToHome)
	r.Get("/org/", reqOrgAdmin, hs.Index)
	r.Get("/org/new", reqGrafinsightAdmin, hs.Index)
	r.Get("/datasources/", reqOrgAdmin, hs.Index)
	r.Get("/datasources/new", reqOrgAdmin, hs.Index)
	r.Get("/datasources/edit/*", reqOrgAdmin, hs.Index)
	r.Get("/org/users", reqOrgAdmin, hs.Index)
	r.Get("/org/users/new", reqOrgAdmin, hs.Index)
	r.Get("/org/users/invite", reqOrgAdmin, hs.Index)
	r.Get("/org/teams", reqCanAccessTeams, hs.Index)
	r.Get("/org/teams/*", reqCanAccessTeams, hs.Index)
	r.Get("/org/apikeys/", reqOrgAdmin, hs.Index)
	r.Get("/dashboard/import/", reqSignedIn, hs.Index)
	r.Get("/configuration", reqGrafinsightAdmin, hs.Index)
	r.Get("/admin", reqGrafinsightAdmin, hs.Index)
	r.Get("/admin/settings", reqGrafinsightAdmin, hs.Index)
	r.Get("/admin/users", reqGrafinsightAdmin, hs.Index)
	r.Get("/admin/users/create", reqGrafinsightAdmin, hs.Index)
	r.Get("/admin/users/edit/:id", reqGrafinsightAdmin, hs.Index)
	r.Get("/admin/orgs", reqGrafinsightAdmin, hs.Index)
	r.Get("/admin/orgs/edit/:id", reqGrafinsightAdmin, hs.Index)
	r.Get("/admin/stats", reqGrafinsightAdmin, hs.Index)
	r.Get("/admin/ldap", reqGrafinsightAdmin, hs.Index)

	r.Get("/styleguide", reqSignedIn, hs.Index)

	r.Get("/plugins", reqSignedIn, hs.Index)
	r.Get("/plugins/:id/", reqSignedIn, hs.Index)
	r.Get("/plugins/:id/edit", reqSignedIn, hs.Index) // deprecated
	r.Get("/plugins/:id/page/:page", reqSignedIn, hs.Index)
	r.Get("/a/:id/*", reqSignedIn, hs.Index) // App Root Page

	r.Get("/d/:uid/:slug", reqSignedIn, redirectFromLegacyPanelEditURL, hs.Index)
	r.Get("/d/:uid", reqSignedIn, redirectFromLegacyPanelEditURL, hs.Index)
	r.Get("/dashboard/db/:slug", reqSignedIn, redirectFromLegacyDashboardURL, hs.Index)
	r.Get("/dashboard/script/*", reqSignedIn, hs.Index)
	r.Get("/dashboard/new", reqSignedIn, hs.Index)
	r.Get("/dashboard-solo/snapshot/*", hs.Index)
	r.Get("/d-solo/:uid/:slug", reqSignedIn, hs.Index)
	r.Get("/d-solo/:uid", reqSignedIn, hs.Index)
	r.Get("/dashboard-solo/db/:slug", reqSignedIn, redirectFromLegacyDashboardSoloURL, hs.Index)
	r.Get("/dashboard-solo/script/*", reqSignedIn, hs.Index)
	r.Get("/import/dashboard", reqSignedIn, hs.Index)
	r.Get("/dashboards/", reqSignedIn, hs.Index)
	r.Get("/dashboards/*", reqSignedIn, hs.Index)
	r.Get("/goto/:uid", reqSignedIn, hs.redirectFromShortURL, hs.Index)

	r.Get("/explore", reqSignedIn, middleware.EnsureEditorOrViewerCanEdit, hs.Index)

	r.Get("/playlists/", reqSignedIn, hs.Index)
	r.Get("/playlists/*", reqSignedIn, hs.Index)
	r.Get("/alerting/", reqEditorRole, hs.Index)
	r.Get("/alerting/*", reqEditorRole, hs.Index)

	// sign up
	r.Get("/verify", hs.Index)
	r.Get("/signup", hs.Index)
	r.Get("/api/user/signup/options", routing.Wrap(GetSignUpOptions))
	r.Post("/api/user/signup", quota("user"), bind(dtos.SignUpForm{}), routing.Wrap(SignUp))
	r.Post("/api/user/signup/step2", bind(dtos.SignUpStep2Form{}), routing.Wrap(hs.SignUpStep2))

	// invited
	r.Get("/api/user/invite/:code", routing.Wrap(GetInviteInfoByCode))
	r.Post("/api/user/invite/complete", bind(dtos.CompleteInviteForm{}), routing.Wrap(hs.CompleteInvite))

	// reset password
	r.Get("/user/password/send-reset-email", hs.Index)
	r.Get("/user/password/reset", hs.Index)

	r.Post("/api/user/password/send-reset-email", bind(dtos.SendResetPasswordEmailForm{}), routing.Wrap(SendResetPasswordEmail))
	r.Post("/api/user/password/reset", bind(dtos.ResetUserPasswordForm{}), routing.Wrap(ResetPassword))

	// dashboard snapshots
	r.Get("/dashboard/snapshot/*", reqNoAuth, hs.Index)
	r.Get("/dashboard/snapshots/", reqSignedIn, hs.Index)

	// api renew session based on cookie
	r.Get("/api/login/ping", quota("session"), routing.Wrap(hs.LoginAPIPing))

	// authed api
	r.Group("/api", func(apiRoute routing.RouteRegister) {
		// user (signed in)
		apiRoute.Group("/user", func(userRoute routing.RouteRegister) {
			userRoute.Get("/", routing.Wrap(GetSignedInUser))
			userRoute.Put("/", bind(models.UpdateUserCommand{}), routing.Wrap(UpdateSignedInUser))
			userRoute.Post("/using/:id", routing.Wrap(UserSetUsingOrg))
			userRoute.Get("/orgs", routing.Wrap(GetSignedInUserOrgList))
			userRoute.Get("/teams", routing.Wrap(GetSignedInUserTeamList))

			userRoute.Post("/stars/dashboard/:id", routing.Wrap(StarDashboard))
			userRoute.Delete("/stars/dashboard/:id", routing.Wrap(UnstarDashboard))

			userRoute.Put("/password", bind(models.ChangeUserPasswordCommand{}), routing.Wrap(ChangeUserPassword))
			userRoute.Get("/quotas", routing.Wrap(GetUserQuotas))
			userRoute.Put("/helpflags/:id", routing.Wrap(SetHelpFlag))
			// For dev purpose
			userRoute.Get("/helpflags/clear", routing.Wrap(ClearHelpFlags))

			userRoute.Get("/preferences", routing.Wrap(GetUserPreferences))
			userRoute.Put("/preferences", bind(dtos.UpdatePrefsCmd{}), routing.Wrap(UpdateUserPreferences))

			userRoute.Get("/auth-tokens", routing.Wrap(hs.GetUserAuthTokens))
			userRoute.Post("/revoke-auth-token", bind(models.RevokeAuthTokenCmd{}), routing.Wrap(hs.RevokeUserAuthToken))
		}, reqSignedInNoAnonymous)

		// users (admin permission required)
		apiRoute.Group("/users", func(usersRoute routing.RouteRegister) {
			usersRoute.Get("/", routing.Wrap(SearchUsers))
			usersRoute.Get("/search", routing.Wrap(SearchUsersWithPaging))
			usersRoute.Get("/:id", routing.Wrap(GetUserByID))
			usersRoute.Get("/:id/teams", routing.Wrap(GetUserTeams))
			usersRoute.Get("/:id/orgs", routing.Wrap(GetUserOrgList))
			// query parameters /users/lookup?loginOrEmail=admin@example.com
			usersRoute.Get("/lookup", routing.Wrap(GetUserByLoginOrEmail))
			usersRoute.Put("/:id", bind(models.UpdateUserCommand{}), routing.Wrap(UpdateUser))
			usersRoute.Post("/:id/using/:orgId", routing.Wrap(UpdateUserActiveOrg))
		}, reqGrafinsightAdmin)

		// team (admin permission required)
		apiRoute.Group("/teams", func(teamsRoute routing.RouteRegister) {
			teamsRoute.Post("/", bind(models.CreateTeamCommand{}), routing.Wrap(hs.CreateTeam))
			teamsRoute.Put("/:teamId", bind(models.UpdateTeamCommand{}), routing.Wrap(hs.UpdateTeam))
			teamsRoute.Delete("/:teamId", routing.Wrap(hs.DeleteTeamByID))
			teamsRoute.Get("/:teamId/members", routing.Wrap(hs.GetTeamMembers))
			teamsRoute.Post("/:teamId/members", bind(models.AddTeamMemberCommand{}), routing.Wrap(hs.AddTeamMember))
			teamsRoute.Put("/:teamId/members/:userId", bind(models.UpdateTeamMemberCommand{}), routing.Wrap(hs.UpdateTeamMember))
			teamsRoute.Delete("/:teamId/members/:userId", routing.Wrap(hs.RemoveTeamMember))
			teamsRoute.Get("/:teamId/preferences", routing.Wrap(hs.GetTeamPreferences))
			teamsRoute.Put("/:teamId/preferences", bind(dtos.UpdatePrefsCmd{}), routing.Wrap(hs.UpdateTeamPreferences))
		}, reqCanAccessTeams)

		// team without requirement of user to be org admin
		apiRoute.Group("/teams", func(teamsRoute routing.RouteRegister) {
			teamsRoute.Get("/:teamId", routing.Wrap(hs.GetTeamByID))
			teamsRoute.Get("/search", routing.Wrap(hs.SearchTeams))
		})

		// org information available to all users.
		apiRoute.Group("/org", func(orgRoute routing.RouteRegister) {
			orgRoute.Get("/", routing.Wrap(GetOrgCurrent))
			orgRoute.Get("/quotas", routing.Wrap(GetOrgQuotas))
		})

		// current org
		apiRoute.Group("/org", func(orgRoute routing.RouteRegister) {
			orgRoute.Put("/", bind(dtos.UpdateOrgForm{}), routing.Wrap(UpdateOrgCurrent))
			orgRoute.Put("/address", bind(dtos.UpdateOrgAddressForm{}), routing.Wrap(UpdateOrgAddressCurrent))
			orgRoute.Get("/users", routing.Wrap(hs.GetOrgUsersForCurrentOrg))
			orgRoute.Post("/users", quota("user"), bind(models.AddOrgUserCommand{}), routing.Wrap(AddOrgUserToCurrentOrg))
			orgRoute.Patch("/users/:userId", bind(models.UpdateOrgUserCommand{}), routing.Wrap(UpdateOrgUserForCurrentOrg))
			orgRoute.Delete("/users/:userId", routing.Wrap(RemoveOrgUserForCurrentOrg))

			// invites
			orgRoute.Get("/invites", routing.Wrap(GetPendingOrgInvites))
			orgRoute.Post("/invites", quota("user"), bind(dtos.AddInviteForm{}), routing.Wrap(AddOrgInvite))
			orgRoute.Patch("/invites/:code/revoke", routing.Wrap(RevokeInvite))

			// prefs
			orgRoute.Get("/preferences", routing.Wrap(GetOrgPreferences))
			orgRoute.Put("/preferences", bind(dtos.UpdatePrefsCmd{}), routing.Wrap(UpdateOrgPreferences))
		}, reqOrgAdmin)

		// current org without requirement of user to be org admin
		apiRoute.Group("/org", func(orgRoute routing.RouteRegister) {
			orgRoute.Get("/users/lookup", routing.Wrap(hs.GetOrgUsersForCurrentOrgLookup))
		})

		// create new org
		apiRoute.Post("/orgs", quota("org"), bind(models.CreateOrgCommand{}), routing.Wrap(CreateOrg))

		// search all orgs
		apiRoute.Get("/orgs", reqGrafinsightAdmin, routing.Wrap(SearchOrgs))

		// orgs (admin routes)
		apiRoute.Group("/orgs/:orgId", func(orgsRoute routing.RouteRegister) {
			orgsRoute.Get("/", routing.Wrap(GetOrgByID))
			orgsRoute.Put("/", bind(dtos.UpdateOrgForm{}), routing.Wrap(UpdateOrg))
			orgsRoute.Put("/address", bind(dtos.UpdateOrgAddressForm{}), routing.Wrap(UpdateOrgAddress))
			orgsRoute.Delete("/", routing.Wrap(DeleteOrgByID))
			orgsRoute.Get("/users", routing.Wrap(hs.GetOrgUsers))
			orgsRoute.Post("/users", bind(models.AddOrgUserCommand{}), routing.Wrap(AddOrgUser))
			orgsRoute.Patch("/users/:userId", bind(models.UpdateOrgUserCommand{}), routing.Wrap(UpdateOrgUser))
			orgsRoute.Delete("/users/:userId", routing.Wrap(RemoveOrgUser))
			orgsRoute.Get("/quotas", routing.Wrap(GetOrgQuotas))
			orgsRoute.Put("/quotas/:target", bind(models.UpdateOrgQuotaCmd{}), routing.Wrap(UpdateOrgQuota))
		}, reqGrafinsightAdmin)

		// orgs (admin routes)
		apiRoute.Group("/orgs/name/:name", func(orgsRoute routing.RouteRegister) {
			orgsRoute.Get("/", routing.Wrap(hs.GetOrgByName))
		}, reqGrafinsightAdmin)

		// auth api keys
		apiRoute.Group("/auth/keys", func(keysRoute routing.RouteRegister) {
			keysRoute.Get("/", routing.Wrap(GetAPIKeys))
			keysRoute.Post("/", quota("api_key"), bind(models.AddApiKeyCommand{}), routing.Wrap(hs.AddAPIKey))
			keysRoute.Delete("/:id", routing.Wrap(DeleteAPIKey))
		}, reqOrgAdmin)

		// Preferences
		apiRoute.Group("/preferences", func(prefRoute routing.RouteRegister) {
			prefRoute.Post("/set-home-dash", bind(models.SavePreferencesCommand{}), routing.Wrap(SetHomeDashboard))
		})

		// Data sources
		apiRoute.Group("/datasources", func(datasourceRoute routing.RouteRegister) {
			datasourceRoute.Get("/", routing.Wrap(hs.GetDataSources))
			datasourceRoute.Post("/", quota("data_source"), bind(models.AddDataSourceCommand{}), routing.Wrap(AddDataSource))
			datasourceRoute.Put("/:id", bind(models.UpdateDataSourceCommand{}), routing.Wrap(UpdateDataSource))
			datasourceRoute.Delete("/:id", routing.Wrap(DeleteDataSourceById))
			datasourceRoute.Delete("/uid/:uid", routing.Wrap(DeleteDataSourceByUID))
			datasourceRoute.Delete("/name/:name", routing.Wrap(DeleteDataSourceByName))
			datasourceRoute.Get("/:id", routing.Wrap(GetDataSourceById))
			datasourceRoute.Get("/uid/:uid", routing.Wrap(GetDataSourceByUID))
			datasourceRoute.Get("/name/:name", routing.Wrap(GetDataSourceByName))
		}, reqOrgAdmin)

		apiRoute.Get("/datasources/id/:name", routing.Wrap(GetDataSourceIdByName), reqSignedIn)

		apiRoute.Get("/plugins", routing.Wrap(hs.GetPluginList))
		apiRoute.Get("/plugins/:pluginId/settings", routing.Wrap(GetPluginSettingByID))
		apiRoute.Get("/plugins/:pluginId/markdown/:name", routing.Wrap(GetPluginMarkdown))
		apiRoute.Get("/plugins/:pluginId/health", routing.Wrap(hs.CheckHealth))
		apiRoute.Any("/plugins/:pluginId/resources", hs.CallResource)
		apiRoute.Any("/plugins/:pluginId/resources/*", hs.CallResource)
		apiRoute.Any("/plugins/errors", routing.Wrap(hs.GetPluginErrorsList))

		apiRoute.Group("/plugins", func(pluginRoute routing.RouteRegister) {
			pluginRoute.Get("/:pluginId/dashboards/", routing.Wrap(GetPluginDashboards))
			pluginRoute.Post("/:pluginId/settings", bind(models.UpdatePluginSettingCmd{}), routing.Wrap(UpdatePluginSetting))
			pluginRoute.Get("/:pluginId/metrics", routing.Wrap(hs.CollectPluginMetrics))
		}, reqOrgAdmin)

		apiRoute.Get("/frontend/settings/", hs.GetFrontendSettings)
		apiRoute.Any("/datasources/proxy/:id/*", reqSignedIn, hs.ProxyDataSourceRequest)
		apiRoute.Any("/datasources/proxy/:id", reqSignedIn, hs.ProxyDataSourceRequest)
		apiRoute.Any("/datasources/:id/resources", hs.CallDatasourceResource)
		apiRoute.Any("/datasources/:id/resources/*", hs.CallDatasourceResource)
		apiRoute.Any("/datasources/:id/health", routing.Wrap(hs.CheckDatasourceHealth))

		// Folders
		apiRoute.Group("/folders", func(folderRoute routing.RouteRegister) {
			folderRoute.Get("/", routing.Wrap(GetFolders))
			folderRoute.Get("/id/:id", routing.Wrap(GetFolderByID))
			folderRoute.Post("/", bind(models.CreateFolderCommand{}), routing.Wrap(hs.CreateFolder))

			folderRoute.Group("/:uid", func(folderUidRoute routing.RouteRegister) {
				folderUidRoute.Get("/", routing.Wrap(GetFolderByUID))
				folderUidRoute.Put("/", bind(models.UpdateFolderCommand{}), routing.Wrap(UpdateFolder))
				folderUidRoute.Delete("/", routing.Wrap(hs.DeleteFolder))

				folderUidRoute.Group("/permissions", func(folderPermissionRoute routing.RouteRegister) {
					folderPermissionRoute.Get("/", routing.Wrap(hs.GetFolderPermissionList))
					folderPermissionRoute.Post("/", bind(dtos.UpdateDashboardAclCommand{}), routing.Wrap(hs.UpdateFolderPermissions))
				})
			})
		})

		// Dashboard
		apiRoute.Group("/dashboards", func(dashboardRoute routing.RouteRegister) {
			dashboardRoute.Get("/uid/:uid", routing.Wrap(hs.GetDashboard))
			dashboardRoute.Delete("/uid/:uid", routing.Wrap(hs.DeleteDashboardByUID))

			dashboardRoute.Get("/db/:slug", routing.Wrap(hs.GetDashboard))
			dashboardRoute.Delete("/db/:slug", routing.Wrap(hs.DeleteDashboardBySlug))

			dashboardRoute.Post("/calculate-diff", bind(dtos.CalculateDiffOptions{}), routing.Wrap(CalculateDashboardDiff))

			dashboardRoute.Post("/db", bind(models.SaveDashboardCommand{}), routing.Wrap(hs.PostDashboard))
			dashboardRoute.Get("/home", routing.Wrap(hs.GetHomeDashboard))
			dashboardRoute.Get("/tags", GetDashboardTags)
			dashboardRoute.Post("/import", bind(dtos.ImportDashboardCommand{}), routing.Wrap(ImportDashboard))

			dashboardRoute.Group("/id/:dashboardId", func(dashIdRoute routing.RouteRegister) {
				dashIdRoute.Get("/versions", routing.Wrap(GetDashboardVersions))
				dashIdRoute.Get("/versions/:id", routing.Wrap(GetDashboardVersion))
				dashIdRoute.Post("/restore", bind(dtos.RestoreDashboardVersionCommand{}), routing.Wrap(hs.RestoreDashboardVersion))

				dashIdRoute.Group("/permissions", func(dashboardPermissionRoute routing.RouteRegister) {
					dashboardPermissionRoute.Get("/", routing.Wrap(hs.GetDashboardPermissionList))
					dashboardPermissionRoute.Post("/", bind(dtos.UpdateDashboardAclCommand{}), routing.Wrap(hs.UpdateDashboardPermissions))
				})
			})
		})

		// Dashboard snapshots
		apiRoute.Group("/dashboard/snapshots", func(dashboardRoute routing.RouteRegister) {
			dashboardRoute.Get("/", routing.Wrap(SearchDashboardSnapshots))
		})

		// Playlist
		apiRoute.Group("/playlists", func(playlistRoute routing.RouteRegister) {
			playlistRoute.Get("/", routing.Wrap(SearchPlaylists))
			playlistRoute.Get("/:id", ValidateOrgPlaylist, routing.Wrap(GetPlaylist))
			playlistRoute.Get("/:id/items", ValidateOrgPlaylist, routing.Wrap(GetPlaylistItems))
			playlistRoute.Get("/:id/dashboards", ValidateOrgPlaylist, routing.Wrap(GetPlaylistDashboards))
			playlistRoute.Delete("/:id", reqEditorRole, ValidateOrgPlaylist, routing.Wrap(DeletePlaylist))
			playlistRoute.Put("/:id", reqEditorRole, bind(models.UpdatePlaylistCommand{}), ValidateOrgPlaylist, routing.Wrap(UpdatePlaylist))
			playlistRoute.Post("/", reqEditorRole, bind(models.CreatePlaylistCommand{}), routing.Wrap(CreatePlaylist))
		})

		// Search
		apiRoute.Get("/search/sorting", routing.Wrap(hs.ListSortOptions))
		apiRoute.Get("/search/", routing.Wrap(Search))

		// metrics
		apiRoute.Post("/tsdb/query", bind(dtos.MetricRequest{}), routing.Wrap(hs.QueryMetrics))
		apiRoute.Get("/tsdb/testdata/gensql", reqGrafinsightAdmin, routing.Wrap(GenerateSQLTestData))
		apiRoute.Get("/tsdb/testdata/random-walk", routing.Wrap(GetTestDataRandomWalk))

		// DataSource w/ expressions
		apiRoute.Post("/ds/query", bind(dtos.MetricRequest{}), routing.Wrap(hs.QueryMetricsV2))

		apiRoute.Group("/alerts", func(alertsRoute routing.RouteRegister) {
			alertsRoute.Post("/test", bind(dtos.AlertTestCommand{}), routing.Wrap(AlertTest))
			alertsRoute.Post("/:alertId/pause", reqEditorRole, bind(dtos.PauseAlertCommand{}), routing.Wrap(PauseAlert))
			alertsRoute.Get("/:alertId", ValidateOrgAlert, routing.Wrap(GetAlert))
			alertsRoute.Get("/", routing.Wrap(GetAlerts))
			alertsRoute.Get("/states-for-dashboard", routing.Wrap(GetAlertStatesForDashboard))
		})

		apiRoute.Get("/alert-notifiers", reqEditorRole, routing.Wrap(GetAlertNotifiers))

		apiRoute.Group("/alert-notifications", func(alertNotifications routing.RouteRegister) {
			alertNotifications.Get("/", routing.Wrap(GetAlertNotifications))
			alertNotifications.Post("/test", bind(dtos.NotificationTestCommand{}), routing.Wrap(NotificationTest))
			alertNotifications.Post("/", bind(models.CreateAlertNotificationCommand{}), routing.Wrap(CreateAlertNotification))
			alertNotifications.Put("/:notificationId", bind(models.UpdateAlertNotificationCommand{}), routing.Wrap(UpdateAlertNotification))
			alertNotifications.Get("/:notificationId", routing.Wrap(GetAlertNotificationByID))
			alertNotifications.Delete("/:notificationId", routing.Wrap(DeleteAlertNotification))
			alertNotifications.Get("/uid/:uid", routing.Wrap(GetAlertNotificationByUID))
			alertNotifications.Put("/uid/:uid", bind(models.UpdateAlertNotificationWithUidCommand{}), routing.Wrap(UpdateAlertNotificationByUID))
			alertNotifications.Delete("/uid/:uid", routing.Wrap(DeleteAlertNotificationByUID))
		}, reqEditorRole)

		// alert notifications without requirement of user to be org editor
		apiRoute.Group("/alert-notifications", func(orgRoute routing.RouteRegister) {
			orgRoute.Get("/lookup", routing.Wrap(GetAlertNotificationLookup))
		})

		apiRoute.Get("/annotations", routing.Wrap(GetAnnotations))
		apiRoute.Post("/annotations/mass-delete", reqOrgAdmin, bind(dtos.DeleteAnnotationsCmd{}), routing.Wrap(DeleteAnnotations))

		apiRoute.Group("/annotations", func(annotationsRoute routing.RouteRegister) {
			annotationsRoute.Post("/", bind(dtos.PostAnnotationsCmd{}), routing.Wrap(PostAnnotation))
			annotationsRoute.Delete("/:annotationId", routing.Wrap(DeleteAnnotationByID))
			annotationsRoute.Put("/:annotationId", bind(dtos.UpdateAnnotationsCmd{}), routing.Wrap(UpdateAnnotation))
			annotationsRoute.Patch("/:annotationId", bind(dtos.PatchAnnotationsCmd{}), routing.Wrap(PatchAnnotation))
			annotationsRoute.Post("/graphite", reqEditorRole, bind(dtos.PostGraphiteAnnotationsCmd{}), routing.Wrap(PostGraphiteAnnotation))
		})

		// short urls
		apiRoute.Post("/short-urls", bind(dtos.CreateShortURLCmd{}), routing.Wrap(hs.createShortURL))
	}, reqSignedIn)

	// admin api
	r.Group("/api/admin", func(adminRoute routing.RouteRegister) {
		adminRoute.Get("/settings", routing.Wrap(AdminGetSettings))
		adminRoute.Post("/users", bind(dtos.AdminCreateUserForm{}), routing.Wrap(AdminCreateUser))
		adminRoute.Put("/users/:id/password", bind(dtos.AdminUpdateUserPasswordForm{}), routing.Wrap(AdminUpdateUserPassword))
		adminRoute.Put("/users/:id/permissions", bind(dtos.AdminUpdateUserPermissionsForm{}), routing.Wrap(AdminUpdateUserPermissions))
		adminRoute.Delete("/users/:id", routing.Wrap(AdminDeleteUser))
		adminRoute.Post("/users/:id/disable", routing.Wrap(hs.AdminDisableUser))
		adminRoute.Post("/users/:id/enable", routing.Wrap(AdminEnableUser))
		adminRoute.Get("/users/:id/quotas", routing.Wrap(GetUserQuotas))
		adminRoute.Put("/users/:id/quotas/:target", bind(models.UpdateUserQuotaCmd{}), routing.Wrap(UpdateUserQuota))
		adminRoute.Get("/stats", routing.Wrap(AdminGetStats))
		adminRoute.Post("/pause-all-alerts", bind(dtos.PauseAllAlertsCommand{}), routing.Wrap(PauseAllAlerts))

		adminRoute.Post("/users/:id/logout", routing.Wrap(hs.AdminLogoutUser))
		adminRoute.Get("/users/:id/auth-tokens", routing.Wrap(hs.AdminGetUserAuthTokens))
		adminRoute.Post("/users/:id/revoke-auth-token", bind(models.RevokeAuthTokenCmd{}), routing.Wrap(hs.AdminRevokeUserAuthToken))

		adminRoute.Post("/provisioning/dashboards/reload", routing.Wrap(hs.AdminProvisioningReloadDashboards))
		adminRoute.Post("/provisioning/plugins/reload", routing.Wrap(hs.AdminProvisioningReloadPlugins))
		adminRoute.Post("/provisioning/datasources/reload", routing.Wrap(hs.AdminProvisioningReloadDatasources))
		adminRoute.Post("/provisioning/notifications/reload", routing.Wrap(hs.AdminProvisioningReloadNotifications))
		adminRoute.Post("/ldap/reload", routing.Wrap(hs.ReloadLDAPCfg))
		adminRoute.Post("/ldap/sync/:id", routing.Wrap(hs.PostSyncUserWithLDAP))
		adminRoute.Get("/ldap/:username", routing.Wrap(hs.GetUserFromLDAP))
		adminRoute.Get("/ldap/status", routing.Wrap(hs.GetLDAPStatus))
	}, reqGrafinsightAdmin)

	// rendering
	r.Get("/render/*", reqSignedIn, hs.RenderToPng)

	// grafinsight.net proxy
	r.Any("/api/gnet/*", reqSignedIn, ProxyGnetRequest)

	// Gravatar service.
	avatarCacheServer := avatar.NewCacheServer()
	r.Get("/avatar/:hash", avatarCacheServer.Handler)

	// Snapshots
	r.Post("/api/snapshots/", reqSnapshotPublicModeOrSignedIn, bind(models.CreateDashboardSnapshotCommand{}), CreateDashboardSnapshot)
	r.Get("/api/snapshot/shared-options/", reqSignedIn, GetSharingOptions)
	r.Get("/api/snapshots/:key", routing.Wrap(GetDashboardSnapshot))
	r.Get("/api/snapshots-delete/:deleteKey", reqSnapshotPublicModeOrSignedIn, routing.Wrap(DeleteDashboardSnapshotByDeleteKey))
	r.Delete("/api/snapshots/:key", reqEditorRole, routing.Wrap(DeleteDashboardSnapshot))

	// Frontend logs
	sourceMapStore := frontendlogging.NewSourceMapStore(hs.Cfg, frontendlogging.ReadSourceMapFromFS)
	r.Post("/log", middleware.RateLimit(hs.Cfg.Sentry.EndpointRPS, hs.Cfg.Sentry.EndpointBurst, time.Now), bind(frontendlogging.FrontendSentryEvent{}), routing.Wrap(NewFrontendLogMessageHandler(sourceMapStore)))
}
