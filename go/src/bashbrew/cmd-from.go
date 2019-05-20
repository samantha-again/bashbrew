package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

func cmdFrom(c *cli.Context) error {
	repos, err := repos(c.Bool("all"), c.Args()...)
	if err != nil {
		return cli.NewMultiError(fmt.Errorf(`failed gathering repo list`), err)
	}

	uniq := c.Bool("uniq")
	namespace := ""
	applyConstraints := c.Bool("apply-constraints")

	for _, repo := range repos {
		r, err := fetch(repo)
		if err != nil {
			return cli.NewMultiError(fmt.Errorf(`failed fetching repo %q`, repo), err)
		}

		for _, entry := range r.Entries() {
			if applyConstraints && r.SkipConstraints(entry) {
				continue
			}

			froms, err := r.DockerFroms(entry)
			if err != nil {
				return cli.NewMultiError(fmt.Errorf(`failed fetching/scraping FROM for %q (tags %q)`, r.RepoName, entry.TagsString()), err)
			}

			fromsString := strings.Join(froms, " ")
			for _, tag := range r.Tags(namespace, uniq, entry) {
				fmt.Printf("%s: %s\n", tag, fromsString)
			}
		}
	}

	return nil
}
