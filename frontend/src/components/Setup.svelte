<script lang="ts">
  import type { MouseEventHandler } from "svelte/elements";
  import { onMount } from "svelte";
  import Logs from "./Logs.svelte";
  import { router } from "@/lib/router.svelte";
  import { svg, RELOAD, TARGET } from "@/lib/svg";
  import {
    NewADBClient,
    SelectDevice,
    KillServer,
    DownloadADB,
    GetDeviceList,
    ConnectToServer,
    SelectDownloadDir,
    SelectADBExecutable,
    GetDefaultSettings,
  } from "@wails/go/main/App";

  let port = $state(5037);
  let paths = $state({ downloadDir: "loading...", adb: "loading..." });
  let selectedDevice = $state<number | null>();
  let devices = $state<string[]>([]);

  onMount(async () => {
    paths = await GetDefaultSettings();
  });

  const onADBSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    if (paths.adb.startsWith("loading")) return;

    const formData = new FormData(e.currentTarget as HTMLFormElement);
    switch (formData.get("action")) {
      case "kill": {
        await KillServer(paths.adb, port);
        devices = [];
        selectedDevice = null;
        return;
      }
      case "start": {
        await NewADBClient(paths.adb, port);
        break;
      }
      case "connect": {
        await ConnectToServer(port);
        break;
      }
      default:
        throw new Error("unkown action");
    }

    await refreshDevices();
  };

  const selectDevice = async (e: SubmitEvent) => {
    e.preventDefault();
    if (selectedDevice == null) return;
    await SelectDevice(selectedDevice, paths.downloadDir);
    router.current = "explore";
  };

  const downloadADB = async () => {
    paths.adb = await DownloadADB();
  };

  const refreshDevices = async () => {
    devices = await GetDeviceList();
    if (devices == null || devices.length === 0) return;
    selectedDevice = 0;
  };

  const selectADBExecutable = async () => {
    paths.adb = await SelectADBExecutable();
  };

  const selectDownloadDir = async () => {
    paths.downloadDir = await SelectDownloadDir();
  };
</script>

<form onsubmit={onADBSubmit}>
  <label>
    <span>Port</span>
    <input required bind:value={port} type="number" />
  </label>

  {@render pathInput(paths, "adb", "ADB executable path", selectADBExecutable)}

  <div>
    <button name="action" value="connect" type="submit">Connect</button>
    <button name="action" value="start" type="submit">Start server</button>
    <button name="action" value="kill" type="submit">Kill server</button>
  </div>
</form>

<form onsubmit={selectDevice}>
  {@render pathInput(paths, "downloadDir", "Download directory path", selectDownloadDir)}

  <div>
    <label>
      <span>Device</span>
      <select required bind:value={selectedDevice}>
        {#if devices == null || devices.length === 0}
          <option>No device</option>
        {/if}

        {#each devices as device, i}
          <option value={i}>{device}</option>
        {/each}
      </select>
    </label>

    <button type="button" title="refresh" onclick={refreshDevices}>
      {@render svg({ d: RELOAD })}
    </button>
  </div>

  <button type="submit" disabled={typeof selectedDevice !== "number"}>Select</button>
</form>

<section>
  <!-- TODO: download progress bar -->
  <!-- NOTE: use resp.ContentLength with io.TeeReader -->
  <!-- TODO: check server is running on given port or any onsubmit errors -->
  <button onclick={downloadADB}>Download ADB</button>
  <Logs />
</section>

{#snippet pathInput(
  values: typeof paths,
  key: keyof typeof paths,
  label: string,
  onclick: MouseEventHandler<HTMLButtonElement>,
)}
  <div>
    <label>
      <span>{label}</span>
      <input required type="text" bind:value={values[key]} />
    </label>
    <button
      {onclick}
      type="button"
      title={`select ${key === "adb" ? "adb executable" : "download directory"}`}>
      {@render svg({ d: TARGET })}
    </button>
  </div>
{/snippet}

<style>
</style>
