// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// To use this file, you need to enable DOM and ES2015.Promise in
// tsconfig.json
//
// "lib": ["dom", "es2015.promise"]

interface JsonResp {
    data?: any;
    error?: any;
}

function grab<T>(uri: Request | string, init?: RequestInit): Promise<T> {
    return fetch(uri, init)
        .then((resp: Response) => {
            return resp.json();
        })
        .then((data: any) => {
            if (data.error) {
                throw new Error(data.error);
            }

            return <T>data.data;
        });
}
