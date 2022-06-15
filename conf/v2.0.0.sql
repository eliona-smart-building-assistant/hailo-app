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

-- Remove constraint to decouple config from core eliona database
alter table hailo.config drop constraint if exists hailo_config_asset_id_fkey;

-- Create table that defines the eliona projects for using the configuration. Within these projects
-- tha app create or update the assets automatically.
alter table hailo.config add column if not exists proj_ids text[];

-- Create table to pair eliona assets with hailo fds devices
-- Note, that existing installations have to migrate the mapping
-- from public.asset.device_pkey -> public.asset.asset_id
create table if not exists hailo.asset
(
    config_id       bigint not null,
    device_id       text not null,
    proj_id         text not null,
    asset_id        integer not null,
    primary key (config_id, device_id, proj_id, asset_id)
);
