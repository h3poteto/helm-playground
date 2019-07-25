package deploy

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"text/tabwriter"

	"github.com/gosuri/uitable"
	"github.com/gosuri/uitable/util/strutil"
	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/timeconv"
)

type Deploy struct {
	client      *helm.Client
	DryRun      bool
	kubeContext string
	kubeConfig  string
}

func New(kubeContext, kubeConfig string) (*Deploy, error) {
	cli, err := NewClient(kubeContext, kubeConfig)
	if err != nil {
		return nil, err
	}
	c := &Deploy{
		client:      cli,
		kubeContext: kubeContext,
		kubeConfig:  kubeConfig,
		DryRun:      true,
	}
	return c, nil
}

func (d *Deploy) Version() (string, error) {
	version, err := d.client.GetVersion()
	if err != nil {
		return "", err
	}
	return version.Version.GetSemVer(), nil
}

func (d *Deploy) NewRelease(chartPath string) (*release.Release, error) {
	chartRequested, err := chartutil.Load(chartPath)
	if err != nil {
		return nil, err
	}

	namespace, _, err := kube.GetConfig(d.kubeContext, d.kubeConfig).Namespace()
	if err != nil {
		return nil, err
	}

	values, err := d.overrideValues()
	if err != nil {
		return nil, err
	}

	res, err := d.client.InstallReleaseFromChart(
		chartRequested,
		namespace,
		helm.ValueOverrides(values),
		helm.ReleaseName("akira"),
		helm.InstallDryRun(d.DryRun),
		helm.InstallReuseName(false),
		helm.InstallDisableHooks(false),
		helm.InstallDisableCRDHook(false),
		helm.InstallSubNotes(false),
		helm.InstallTimeout(300),
		helm.InstallWait(false),
		helm.InstallDescription(""),
	)
	if err != nil {
		return nil, err
	}
	release := res.GetRelease()
	if release == nil {
		return nil, nil
	}
	return release, nil
}

func (d *Deploy) overrideValues() ([]byte, error) {
	base := map[string]interface{}{
		"namespace":         os.Getenv("NAMESPACE"),
		"ses_smtp_user":     os.Getenv("SES_SMTP_USER"),
		"ses_smtp_password": os.Getenv("SES_SMTP_PASSWORD"),
		"image": map[string]interface{}{
			"tag": os.Getenv("IMAGE_TAG"),
		},
	}

	return yaml.Marshal(base)
}

func (d *Deploy) PrintRelease(rel *release.Release) error {
	if rel == nil {
		return nil
	}

	fmt.Printf("NAME:    %s\n", rel.Name)
	if !d.DryRun {
		status, err := d.client.ReleaseStatus(rel.Name)
		if err != nil {
			return err
		}
		printStatus(os.Stdout, status)
	}
	return nil
}

func printStatus(out io.Writer, res *services.GetReleaseStatusResponse) {
	if res.Info.LastDeployed != nil {
		fmt.Fprintf(out, "LAST DEPLOYED: %s\n", timeconv.String(res.Info.LastDeployed))
	}
	fmt.Fprintf(out, "NAMESPACE: %s\n", res.Namespace)
	fmt.Fprintf(out, "STATUS: %s\n", res.Info.Status.Code)
	fmt.Fprintf(out, "\n")
	if len(res.Info.Status.Resources) > 0 {
		re := regexp.MustCompile("  +")

		w := tabwriter.NewWriter(out, 0, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "RESOURCES:\n%s\n", re.ReplaceAllString(res.Info.Status.Resources, "\t"))
		w.Flush()
	}
	if res.Info.Status.LastTestSuiteRun != nil {
		lastRun := res.Info.Status.LastTestSuiteRun
		fmt.Fprintf(out, "TEST SUITE:\n%s\n%s\n\n%s\n",
			fmt.Sprintf("Last Started: %s", timeconv.String(lastRun.StartedAt)),
			fmt.Sprintf("Last Completed: %s", timeconv.String(lastRun.CompletedAt)),
			formatTestResults(lastRun.Results))
	}

	if len(res.Info.Status.Notes) > 0 {
		fmt.Fprintf(out, "NOTES:\n%s\n", res.Info.Status.Notes)
	}
}

func formatTestResults(results []*release.TestRun) string {
	tbl := uitable.New()
	tbl.MaxColWidth = 50
	tbl.AddRow("TEST", "STATUS", "INFO", "STARTED", "COMPLETED")
	for i := 0; i < len(results); i++ {
		r := results[i]
		n := r.Name
		s := strutil.PadRight(r.Status.String(), 10, ' ')
		i := r.Info
		ts := timeconv.String(r.StartedAt)
		tc := timeconv.String(r.CompletedAt)
		tbl.AddRow(n, s, i, ts, tc)
	}
	return tbl.String()
}
