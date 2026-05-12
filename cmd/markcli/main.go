/*
Copyright © 2026 Andrii Olkhovych <andrii.olkhovych@tum.de>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"

	"github.com/def4alt/markd/internal/bookmark"
	"github.com/def4alt/markd/internal/cli"
	"github.com/def4alt/markd/internal/platform/clock"
	"github.com/def4alt/markd/internal/platform/ids"
	"github.com/def4alt/markd/internal/storage/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "markd.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := sqlite.Migrate(context.Background(), db); err != nil {
		panic(err)
	}

	repo := sqlite.NewBookmarkRepository(db)
	clock := clock.Real{}
	idgen := ids.UUID{}

	svc := bookmark.NewService(repo, clock, idgen)

	root := cli.NewRootCmd(svc)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
