<script lang="ts">
  import type { Attachment } from "svelte/attachments";
  import type { MouseEventHandler } from "svelte/elements";

  let { selected }: { selected: string[] } = $props();

  let menu = $state({
    height: 0,
    width: 0,
    current: "",
    isOpen: false,
  });
  let positions = $state({
    cursor: { x: 0, y: 0 },
    window: { height: 0, width: 0 },
  });

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
    const row = e.target.closest("tbody > tr");
    if (row == null || e.target.closest("td > input") != null) return;
    positions.cursor = { x: e.clientX, y: e.clientY };
    positions.window = { height: window.innerHeight, width: window.innerWidth };

    const { window: win, cursor } = positions;
    if (win.height - cursor.y < menu.height) positions.cursor.y -= menu.height;
    if (win.width - cursor.x < menu.width) positions.cursor.x -= menu.width;

    menu.current = row.querySelector<HTMLInputElement>("td > input")?.value || "";
    menu.isOpen = true;
  };
</script>

<svelte:window {oncontextmenu} {onclick} />

{#if menu.isOpen}
  <menu {@attach getContextMenuDimension}> </menu>
{/if}

<style></style>
