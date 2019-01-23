/*
   Copyright 2019 Artem Sidorenko <artem@posteo.de>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package generate

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

// printProgress prints the current processing progress on given output using the given channels.
// Routine exists then context ctx cancels
func printProgress( // nolint: gocyclo
	ctx context.Context,
	out io.Writer,
	ctagscounter <-chan bool,
	cmaxtags <-chan int,
	cissuescounter <-chan bool,
	cmaxissues <-chan int,
	cmrscounter <-chan bool,
	cmaxmrs <-chan int,
) {

	go func() {
		var tagscounter int
		var issuescounter int
		var mrscounter int
		maxtags := "X"
		maxissues := "X"
		maxmrs := "X"

		for {
			select {
			case <-ctx.Done():
				fmt.Fprintf(out, "\n") // nolint: errcheck
				return
			case _, ok := <-ctagscounter:
				if ok {
					tagscounter++
				}
			case v, ok := <-cmaxtags:
				if ok {
					maxtags = strconv.Itoa(v)
				}
			case _, ok := <-cissuescounter:
				if ok {
					issuescounter++
				}
			case v, ok := <-cmaxissues:
				if ok {
					maxissues = strconv.Itoa(v)
				}
			case _, ok := <-cmrscounter:
				if ok {
					mrscounter++
				}
			case v, ok := <-cmaxmrs:
				if ok {
					maxmrs = strconv.Itoa(v)
				}
			}

			fmt.Fprintf(out, // nolint: errcheck
				"\rProgress: %v/%v tags, %v/%v issues, %v/%v MRs/PRs",
				tagscounter, maxtags,
				issuescounter, maxissues,
				mrscounter, maxmrs,
			)
		}
	}()

}
