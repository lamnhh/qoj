# QHH Online Judge

An online judge that uses problem structure similar to that of Themis

# Requirement

- `go >= 1.11` (any version that support Go Modules)
- `node`
- `yarn`

# Usage

- `yarn install`.
- `yarn build`.
- `go run main.go`

The server will be hosted at `localhost:3000`

For now, uploading problem can only be done using a placeholder page at `/static/test-upload-problem.html`.

# To-Do

- [ ] Admin dashboard (for managing problems, contests, test data).
- [x] Add responsive CSS.
- [ ] Improve judger (`./timeout`).
- [ ] Docker for deployment.

# License

MIT License

Copyright (c) 2020 Lam Nguyen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
