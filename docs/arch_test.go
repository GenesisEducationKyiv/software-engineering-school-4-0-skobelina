package docs

import (
	"testing"

	archgo "github.com/fdaines/arch-go/api"
	config "github.com/fdaines/arch-go/api/configuration"
)

func TestArchitecture(t *testing.T) {
	configuration := config.Config{
		DependenciesRules: []*config.DependenciesRule{
			{
				Package: "api",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"internal/cron-jobs",
						"internal/mails",
						"internal/rates",
						"internal/subscribers",
						"pkg/utils",
						"pkg/repo",
					},
				},
			},
			{
				Package: "internal/constants",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"internal/cron-jobs",
						"internal/mails",
						"internal/rates",
						"internal/subscribers",
						"pkg/utils",
					},
				},
			},
			{
				Package: "docs",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"internal/cron-jobs",
						"internal/mails",
						"internal/rates",
						"internal/subscribers",
					},
				},
			},
			{
				Package: "internal/cron-jobs",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"internal/cron-jobs",
						"pkg/utils",
						"internal/mails",
						"internal/rates",
						"internal/subscribers",
						"internal",
					},
				},
			},
			{
				Package: "internal/mails",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"internal/cron-jobs",
						"pkg/utils",
					},
				},
			},
			{
				Package: "internal/rates",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"internal/cron-jobs",
						"pkg/utils",
					},
				},
			},
			{
				Package: "internal/subscribers",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"internal/cron-jobs",
						"pkg/utils",
						"internal/constants",
						"internal",
					},
				},
			},
			{
				Package: "internal/cron-jobs",
				ShouldNotDependsOn: &config.Dependencies{
					Internal: []string{
						"api",
					},
				},
			},
			{
				Package: "pkg/repo",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"pkg/utils",
					},
				},
			},
			{
				Package: "pkg/utils",
				ShouldNotDependsOn: &config.Dependencies{
					Internal: []string{
						"api",
						"internal/cron-jobs",
					},
				},
			},
		},
	}
	moduleInfo := config.Load("github.com/skobelina/currency_converter")

	result := archgo.CheckArchitecture(moduleInfo, configuration)

	if !result.Passes {
		t.Fatal("Project doesn't pass architecture tests")
	}
}
