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
						"domains",
						"utils",
						"repo",
					},
				},
			},
			{
				Package: "constants",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"domains",
						"utils",
					},
				},
			},
			{
				Package: "docs",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"domains",
					},
				},
			},
			{
				Package: "domains/cron-jobs",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"domains",
						"utils",
						"domains/mails",
						"domains/rates",
						"domains/subscribers",
					},
				},
			},
			{
				Package: "domains/mails",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"domains",
						"utils",
					},
				},
			},
			{
				Package: "domains/rates",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"domains",
						"utils",
					},
				},
			},
			{
				Package: "domains/subscribers",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"domains",
						"utils",
						"constants",
					},
				},
			},
			{
				Package: "domains",
				ShouldNotDependsOn: &config.Dependencies{
					Internal: []string{
						"api",
					},
				},
			},
			{
				Package: "repo",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"utils",
					},
				},
			},
			{
				Package: "utils",
				ShouldNotDependsOn: &config.Dependencies{
					Internal: []string{
						"api",
						"domains",
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
