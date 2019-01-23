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
		cmaxtags chan<- int,
	) (
		ctags <-chan data.Tag,
	)
	Issues(
		ctx context.Context,
		cerr chan<- error,
		cmaxissues chan<- int,
	) (
		cissues <-chan data.Issue,
	)
	MRs(
		ctx context.Context,
		cerr chan<- error,
		cmaxmrs chan<- int,
	) (
		cmr <-chan data.MR,
	)
	GetNewTagURL(string) (string, error)
	RepositoryExists() (bool, error)
}
