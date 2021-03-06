// This file is part of jsonapi
//
// jsonapi is distributed in two licenses: The Mozilla Public License,
// v. 2.0 and the GNU Lesser Public License.
//
// jsonapi is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE.
//
// See LICENSE.txt for further information.

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
