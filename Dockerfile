# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM openpitrix/openpitrix-builder as builder

WORKDIR /go/src/openpitrix.io/runtime-provider-kubernetes/
COPY . .

RUN mkdir -p /openpitrix_bin
RUN CGO_ENABLED=0 GOOS=linux GOBIN=/openpitrix_bin go install -ldflags '-w -s' -tags netgo openpitrix.io/runtime-provider-kubernetes/cmd/...

FROM alpine:3.7
RUN apk add --update ca-certificates && update-ca-certificates
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /openpitrix_bin/runtime-provider /usr/local/bin/

CMD ["sh"]