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

package helpers

// GetChannelValuesInt gets values from given channel and returns them as slice
func GetChannelValuesInt(c <-chan int) []int {
	return toIntSlice(getChannelValues(c))
}

func getChannelValues(c interface{}) []interface{} {
	var r []interface{}
	fAddInts := func(c <-chan int) {
		for i := range c {
			r = append(r, i)
		}
	}

	switch c := c.(type) {
	case <-chan int:
		fAddInts(c)
	default:
		panic("Not supported type")
	}

	return r
}

func toIntSlice(si []interface{}) []int {
	var r []int

	for _, i := range si {
		r = append(r, i.(int))
	}

	return r
}
