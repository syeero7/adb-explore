import { createSubscriber } from "svelte/reactivity";
import { EventsOn, EventsOff } from "@wails/runtime/runtime";
import { getContext, setContext } from "svelte";

const LOGS_KEY = "logs";

class Logs {
  #data;
  #subscribe;

  constructor() {
    this.#data = $state<string[]>([]);
    this.#subscribe = createSubscriber((update) => {
      const eventName = "log-msg";
      const handler = (msg: string) => {
        this.#data.push(msg);
        console.log(msg);
        update();
      };

      EventsOn(eventName, handler);
      return () => EventsOff(eventName);
    });
  }

  get data() {
    this.#subscribe();
    return this.#data;
  }

  clear() {
    this.#data = [];
  }
}

export function initLogs() {
  const logs = new Logs();
  setContext(LOGS_KEY, logs);
  return logs;
}

export function useLogs() {
  const logs = getContext<Logs>(LOGS_KEY);
  return logs;
}
