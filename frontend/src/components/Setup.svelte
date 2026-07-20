<script lang="ts">
  import { NewADBClient, SelectDevice, KillServer, DownloadADB } from "@wails/go/main/App";
  import Logs from "./Logs.svelte";

  let port = $state(5037);
  let adbPath = $state("/usr/bin/adb");
  let selectedDevice = $state<number>();
  let devices = $state<string[]>([]);

  async function startADB(e: SubmitEvent) {
    e.preventDefault();
    devices = await NewADBClient(adbPath, port);
  }

  async function selectDevice(e: SubmitEvent) {
    e.preventDefault();
    if (selectedDevice == null) return;
    await SelectDevice(selectedDevice);
  }

  async function killServer() {
    await KillServer(adbPath, port);
  }

  async function downloadADB() {
    adbPath = await DownloadADB();
  }
</script>

<form onsubmit={startADB}>
  <label>
    <span>Port</span>
    <input required bind:value={port} type="number" />
  </label>

  <label>
    <span>ADB executable path</span>
    <input required type="text" bind:value={adbPath} />
  </label>
  <!-- TODO: try to auto detect adb executable path -->
  <!-- TODO: allow selecting adb execuable using wails select/open file dialog -->

  <button type="submit">Start</button>
</form>

<form onsubmit={selectDevice}>
  <label>
    <span>Device</span>
    <select required bind:value={selectedDevice}>
      <option>No device</option>

      {#each devices as device, i}
        <option value={i}>{device}</option>
      {/each}
    </select>
  </label>

  <button disabled={typeof selectedDevice !== "number"}>Select</button>
</form>

<div>
  <!-- TODO: download progress bar -->
  <!-- NOTE: use resp.ContentLength with io.TeeReader -->

  <!-- TODO: check server is running on given port or any onsubmit errors -->
  <button type="button" onclick={downloadADB}>Download ADB</button>
  <button type="button" onclick={killServer}>Kill ADB server</button>
  <Logs />
</div>

<style>
</style>
