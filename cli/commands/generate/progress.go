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

// printProgress prints the current processing progress on given output
// using the input on returned channels.
// Print goroutine exists then context ctx cancels or all channels are closed
func printProgress( // nolint: gocyclo
	ctx context.Context,
	out io.Writer,
) (
	ctagscounter chan<- bool,
	cmaxtags chan<- int,
	cissuescounter chan<- bool,
	cmaxissues chan<- int,
	cmrscounter chan<- bool,
	cmaxmrs chan<- int,
) {
	lctagscounter := make(chan bool)
	lcmaxtags := make(chan int)
	lcissuescounter := make(chan bool)
	lcmaxissues := make(chan int)
	lmrscounter := make(chan bool)
	lcmaxmrs := make(chan int)

	go func() {
		var tagscounter int
		var issuescounter int
		var mrscounter int
		maxtags := "X"
		maxissues := "X"
		maxmrs := "X"

		// print newline character when leaving the progress printing routine
		defer func() {
			fmt.Fprintf(out, "\n") // nolint: errcheck
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-lctagscounter:
				if ok {
					tagscounter++
				} else {
					lctagscounter = nil
				}
			case v, ok := <-lcmaxtags:
				if ok {
					maxtags = strconv.Itoa(v)
				} else {
					lcmaxtags = nil
				}
			case _, ok := <-lcissuescounter:
				if ok {
					issuescounter++
				} else {
					lcissuescounter = nil
				}
			case v, ok := <-lcmaxissues:
				if ok {
					maxissues = strconv.Itoa(v)
				} else {
					lcmaxissues = nil
				}
			case _, ok := <-lmrscounter:
				if ok {
					mrscounter++
				} else {
					lmrscounter = nil
				}
			case v, ok := <-lcmaxmrs:
				if ok {
					maxmrs = strconv.Itoa(v)
				} else {
					lcmaxmrs = nil
				}
			}

			if lctagscounter == nil && lcmaxtags == nil &&
				lcissuescounter == nil && lcmaxissues == nil && lmrscounter == nil && lcmaxmrs == nil {
				return
			}

			fmt.Fprintf(out, // nolint: errcheck
				"\rProgress: %v/%v tags, %v/%v issues, %v/%v MRs/PRs",
				tagscounter, maxtags,
				issuescounter, maxissues,
				mrscounter, maxmrs,
			)
		}
	}()

	return lctagscounter, lcmaxtags,
		lcissuescounter, lcmaxissues,
		lmrscounter, lcmaxmrs
}
