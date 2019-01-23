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
func printProgress(
	ctx context.Context,
	out io.Writer,
	ctagscounter <-chan bool,
	cmaxtags <-chan int,
) {

	go func() {
		var tagscounter int
		maxtags := "X"

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
			}

			fmt.Fprintf(out, "\rProgress: %v/%v tags", tagscounter, maxtags) // nolint: errcheck
		}
	}()

}
