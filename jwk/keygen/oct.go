/*
 * Copyright 2016 Fabrício Godoy
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package keygen

import "crypto/rand"

const (
	// MinimumOCTKeySize defines the minimum recommended key size for symmetric
	// keys.
	MinimumOCTKeySize = 512
)

func newSymKey(alg string, size int) ([]byte, error) {
	if size < MinimumOCTKeySize {
		return nil, TooSmallKeySize{MinimumOCTKeySize, size}
	}

	buf := make([]byte, size/8)
	if _, err := rand.Read(buf); err != nil {
		return nil, ErrorGeneratingKey(err.Error())
	}

	return buf, nil
}
