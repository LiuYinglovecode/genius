package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"code.htres.cn/casicloud/adc-genius/pkg/util/flat"
	"code.htres.cn/casicloud/adc-genius/pkg/util/slice"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"

	"bufio"
	"flag"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/manifest"
	helmchart "k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/renderutil"
	"k8s.io/helm/pkg/timeconv"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	hgetter "k8s.io/helm/pkg/getter"
)

func TestLoadDir(t *testing.T) {
	c, err := chartutil.Load("testdata/frobnitz-1.2.3.tgz")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}

	fmt.Printf("%v", c)
}

func TestK8sClient(t *testing.T) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", clientset)
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func int32Ptr(i int32) *int32 { return &i }

func ginMain() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

var repoURL = "https://kubernetes-charts.storage.googleapis.com/redis-9.3.0.tgz"

func TestRemoteHelm(t *testing.T) {
	u, err := url.Parse(repoURL)
	if err != nil {
		t.Error(err)
	}
	httpgetter, err := hgetter.NewHTTPGetter(u.String(), "", "", "")

	if err != nil {
		t.Error(err)
	}

	data, err := httpgetter.Get(u.String())

	if err != nil {
		t.Error(err)
	}

	r := bytes.NewReader(data.Bytes())

	chart, err := chartutil.LoadArchive(r)

	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, chart)
	assert.Equal(t, chart.Metadata.Name, "redis")
	b, err := json.Marshal(chart.Metadata)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))

	// print readme
	readme := findReadme(chart.Files)
	assert.NotNil(t, readme)
	fmt.Println(string(readme.Value))
	// print values
	vl := chartutil.FromYaml(chart.Values.Raw)
	f, err := flat.Flatten(vl, nil)
	if err != nil {
		t.Error(err)
	}
	v, err := json.Marshal(f)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(v))
	now := timeconv.Now()
	// test Render
	renderOpts := renderutil.Options{
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      "test_release",
			IsInstall: true,
			IsUpgrade: false,
			Time:      now,
			Namespace: "default",
		},
		KubeVersion: "1.11.5",
	}

	assert.NotNil(t, renderOpts)

	config := &helmchart.Config{Raw: string(chart.Values.Raw), Values: map[string]*helmchart.Value{}}
	renderedTemplates, err := renderutil.Render(chart, config, renderOpts)
	if err != nil {
		t.Fatal(err)
	}

	listManifests := manifest.SplitManifests(renderedTemplates)
	for _, manifest := range listManifests {
		fmt.Println(chartutil.ToYaml(manifest))
	}

	assert.NotNil(t, renderedTemplates)
}

var readmeFileNames = []string{"readme.md", "readme.txt", "readme"}

func findReadme(files []*any.Any) (file *any.Any) {
	for _, file := range files {
		if slice.ContainsString(readmeFileNames, strings.ToLower(file.TypeUrl), nil) {
			return file
		}
	}
	return nil
}
