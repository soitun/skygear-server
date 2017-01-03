// Copyright 2015-present Oursky Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package migration

import "github.com/jmoiron/sqlx"

type revision_1981535c8aeb struct {
}

func (r *revision_1981535c8aeb) Version() string {
	return "1981535c8aeb"
}

func (r *revision_1981535c8aeb) Up(tx *sqlx.Tx) error {
	if _, err := tx.Exec(`ALTER TABLE _device ADD COLUMN topic text;`); err != nil {
		return err
	}
	return nil
}

func (r *revision_1981535c8aeb) Down(tx *sqlx.Tx) error {
	if _, err := tx.Exec(`ALTER TABLE _device DROP COLUMN topic;`); err != nil {
		return err
	}
	return nil
}
