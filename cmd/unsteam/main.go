package main

import (
	"fmt"
	"os"
	"strconv"
	"unsteam/internal/pkg"

	gtf "github.com/noneback/go-taskflow"
)

type Context struct {
	App      uint
	Depot    uint
	Manifest pkg.Manifest
}

func main() {
	tf := gtf.NewTaskFlow("main")

	branchMap := map[string]uint{
		"download": 0,
	}
	branch := uint(0)

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if val, ok := branchMap[arg]; ok {
			branch = val
		} else {
			println("invalid branch")
			return
		}
	} else {
		println("need branch")
		return
	}

	ctx := &Context{}

	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-app":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					ctx.App = uint(val)
				}
				i++
			}
		case "-depot":
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					ctx.Depot = uint(val)
				}
				i++
			}
		}
	}

	cond := tf.NewCondition("branch", func() uint {
		return branch
	})

	sf1 := tf.NewSubflow("download-depot-flow", func(sf *gtf.Subflow) {
		t1, t2, t3, t4, t5, t6 :=
			sf.NewTask("fetch-manifest", func() {
				data, err := pkg.RequestJSON[pkg.App](fmt.Sprintf("https://manifest.steam.run/api/depot/%d", ctx.App), nil)
				if err != nil {
					panic(err)
				}

				for _, m := range data.Manifests {
					if fmt.Sprint(m.DepotId) == fmt.Sprint(ctx.Depot) {
						ctx.Manifest = m
						fmt.Println("fetched manifest")
						return
					}
				}

				panic("failed to find manifest")
			}),
			sf.NewTask("fetch-request-code", func() {
				id := ctx.Manifest.Id
				data, err := pkg.RequestJSON[pkg.RequestCode]("https://manifest.steam.run/api/manifest/"+id, nil)
				if err != nil {
					panic(err)
				}

				ctx.Manifest.RequestCode = data.Content
				fmt.Println("fetched request code")
			}),
			sf.NewTask("download-manifest", func() {
				manifest := ctx.Manifest
				data, err := pkg.RequestBytes(fmt.Sprintf("https://steampipe.akamaized.net/depot/%d/manifest/%s/5/%s", manifest.DepotId, manifest.Id, manifest.RequestCode), nil)
				if err != nil {
					panic(err)
				}

				if err = pkg.EnsureDir(".data"); err != nil {
					panic(err)
				}

				if err = pkg.WriteFile(".data/manifest", data); err != nil {
					panic(err)
				}

				if err = pkg.Extract(".data/manifest"); err != nil {
					panic(err)
				}

				fmt.Println("downloaded manifest")
			}),
			sf.NewTask("download-key", func() {
				apiKey, err := pkg.Env("US_API_KEY")
				if err != nil {
					panic(err)
				}

				depotId := fmt.Sprint(ctx.Manifest.DepotId)
				depotKey, err := pkg.RequestJSON[pkg.DepotKey](fmt.Sprintf("https://unsteam.cloudflare-delivery914.workers.dev/key?id=%s", depotId), &pkg.Header{
					Key: "Authorization",
					Value: func() string {
						return fmt.Sprintf("Bearer %s", apiKey)
					},
				})
				if err != nil {
					panic(err)
				}

				data := []byte(depotId + ";" + depotKey.Value)

				if err = pkg.EnsureDir(".data"); err != nil {
					panic(err)
				}

				if err = pkg.WriteFile(".data/k", data); err != nil {
					panic(err)
				}

				fmt.Println("downloaded depot key")
			}),
			sf.NewTask("execute-depotdownloader", func() {
				fmt.Println("executing depotdownloader, please wait")

				if err := pkg.Execute(
					"tools/depotdownloader/linux/DepotDownloaderMod",
					[][]string{
						{"-app", fmt.Sprint(ctx.App)},
						{"-depotkeys", ".data/k"},
						{"-depot", fmt.Sprint(ctx.Manifest.DepotId)},
						{"-manifest", ctx.Manifest.Id},
						{"-manifestfile", ".data/z"},
					},
				); err != nil {
					panic(err)
				}
			}),
			sf.NewTask("cleanup", func() {
				fmt.Println("cleaning up")
				if err := pkg.DeleteDir(".data"); err != nil {
					panic(err)
				}
			})

		t2.Succeed(t1)
		t3.Succeed(t2)
		t4.Succeed(t2)
		t5.Succeed(t4, t3)
		t6.Succeed(t5)
	})

	cond.Precede(sf1)

	exec := gtf.NewExecutor(5)
	exec.Run(tf).Wait()
}
