--  This file is part of the eliona project.
--  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
--  ______ _ _
-- |  ____| (_)
-- | |__  | |_  ___  _ __   __ _
-- |  __| | | |/ _ \| '_ \ / _` |
-- | |____| | | (_) | | | | (_| |
-- |______|_|_|\___/|_| |_|\__,_|
--
--  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
--  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
--  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
--  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
--  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

create schema if not exists hailo;

-- Table for configuring one or more connections to Hailo Digital Hubs.
-- This table should be made editable by eliona frontend.
create table if not exists hailo.config
(
    app_id           bigserial primary key,
    config           json      not null,
    enable           boolean   default false,
    description      text,
    asset_id         integer,
    interval_sec     integer not null,
    auth_timeout     integer not null default 5,
    request_timeout  integer not null default 120,
    inactive_timeout integer,
    active           boolean default false,
    proj_ids         text[]
);

-- Makes the new objects available for all other init steps
commit;