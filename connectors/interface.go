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

package connectors

import (
	"context"

	"github.com/artem-sidorenko/chagen/data"
)

// Connector describes the interface of connectors to the data sources
type Connector interface {
	Tags(
		ctx context.Context,
		cerr chan<- error,
	) (
		ctags <-chan data.Tag,
		ctagsscounter <-chan bool,
		cmaxtags <-chan int,
	)
	Issues(
		ctx context.Context,
		cerr chan<- error,
	) (
		cissues <-chan data.Issue,
		cissuescounter <-chan bool,
		cmaxissues <-chan int,
	)
	MRs(
		ctx context.Context,
		cerr chan<- error,
	) (
		cmr <-chan data.MR,
		cmrscounter <-chan bool,
		cmaxmrs <-chan int,
	)
	GetNewTagURL(string) (string, error)
	RepositoryExists() (bool, error)
}
