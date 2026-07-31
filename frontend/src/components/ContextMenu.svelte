<script lang="ts">
  import type { Attachment } from "svelte/attachments";
  import type { MouseEventHandler } from "svelte/elements";
  import InputDialog, { type Dialog } from "./InputDialog.svelte";
  import { svg, CREATE_DIR, DELETE, DOWNLOAD, OPEN_DIR, RENAME, UPLOAD } from "@/lib/svg";
  import {
    toDir,
    directory,
    isStorageDir,
    download,
    upload,
    deleteEntry,
    basename,
    rename,
  } from "@/lib/fs.svelte";

  let { selected }: { selected: string[] } = $props();

  let menu = $state({ height: 0, width: 0, isOpen: false });
  let currentItem = $state("");
  let positions = $state({ cursor: { x: 0, y: 0 }, window: { height: 0, width: 0 } });
  let dialog = $state<Dialog>({
    isOpen: false,
    value: "",
    title: "",
    oldValue: "",
    onSubmit: async () => {},
  });

  const isFile = $derived(currentItem.startsWith("1|"));

  const getContextMenuDimension: Attachment<HTMLMenuElement> = (node) => {
    menu.height = node.offsetHeight;
    menu.width = node.offsetWidth;
  };

  const onclick: MouseEventHandler<Window> = () => {
    menu.isOpen = false;
  };

  const oncontextmenu: MouseEventHandler<Window> = (e) => {
    e.preventDefault();

    if (!(e.target instanceof HTMLElement)) return;
    const row = e.target.closest<HTMLTableRowElement>("tbody > tr");
    if (row == null || e.target.closest("td > input") != null) return;
    positions.cursor = { x: e.clientX, y: e.clientY };
    positions.window = { height: window.innerHeight, width: window.innerWidth };

    const { window: win, cursor } = positions;
    if (win.height - cursor.y < menu.height) positions.cursor.y -= menu.height;
    if (win.width - cursor.x < menu.width) positions.cursor.x -= menu.width;

    currentItem = row.querySelector<HTMLInputElement>("td > input")?.value || "";
    menu.isOpen = true;
  };

  const renenameDialog = () => {
    dialog.title = `Rename ${isFile ? "File" : "Directory"}`;
    const oldName = basename(currentItem);
    dialog.oldValue = oldName;
    dialog.value = oldName;
    dialog.isOpen = true;
    dialog.onSubmit = async (e: SubmitEvent) => {
      e.preventDefault();
      const newName = dialog.value.trim();
      if (newName == dialog.oldValue) return;
      await rename(directory.current, oldName, newName);
      dialog.isOpen = false;
    };
  };
</script>

<svelte:window {oncontextmenu} {onclick} />

<InputDialog bind:dialog />

{#if menu.isOpen}
  <menu {@attach getContextMenuDimension}>
    {const isStorage = isStorageDir(directory.current)}

    {@render item("open", OPEN_DIR, () => toDir(currentItem), isFile && !isStorage)}
    {@render item("download", DOWNLOAD, download("default", [currentItem]), isStorage)}
    {@render item("download to", DOWNLOAD, download("select", [currentItem]), isStorage)}
    {@render item("delete", DELETE, deleteEntry([currentItem]), isStorage)}
    {@render item("rename", RENAME, renenameDialog, isStorage)}
    <hr />
    {@render item("upload files", UPLOAD, upload("files", directory.current), isStorage)}
    {@render item("upload directory", UPLOAD, upload("dir", directory.current), isStorage)}
    {@render item("create directory", CREATE_DIR, () => {}, isStorage)}
    <hr />

    {const isDisabled = isStorage || selected.length === 0}
    {@render item("download selected", DOWNLOAD, download("default", selected), isDisabled)}
    {@render item("download selected to", DOWNLOAD, download("select", selected), isDisabled)}
    {@render item("delete selected", DELETE, deleteEntry(selected), isDisabled)}
  </menu>
{/if}

{#snippet item(
  name: string,
  icon: string,
  onclick: MouseEventHandler<HTMLButtonElement>,
  disabled = false,
)}
  <li>
    <button {onclick} disabled={disabled!!}>
      {@render svg({ d: icon })}
      <span>{name}</span>
    </button>
  </li>
{/snippet}

<style></style>
