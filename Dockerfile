#  This file is part of the eliona project.
#  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
#  ______ _ _
# |  ____| (_)
# | |__  | |_  ___  _ __   __ _
# |  __| | | |/ _ \| '_ \ / _` |
# | |____| | | (_) | | | | (_| |
# |______|_|_|\___/|_| |_|\__,_|
#
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
#  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
#  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
#  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
#  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

FROM eliona/base-golang:latest-1-alpine AS build

RUN apk add git

WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN DATE=$(date) && \
    GIT_COMMIT=$(git rev-list -1 HEAD) && \
    go build -ldflags "-X 'hailo/apiservices.BuildTimestamp=$DATE' -X 'hailo/apiservices.GitCommit=$GIT_COMMIT'" -o ../app

FROM eliona/base-alpine:latest AS target

COPY --from=build /app ./
COPY conf/*.sql ./conf/
COPY eliona/*.json ./eliona/
COPY apiserver/openapi.json /
COPY metadata.json /

ENV APPNAME=hailo

ENV TZ=Europe/Zurich
CMD [ "/app" ]
