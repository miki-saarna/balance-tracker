import { ref } from "vue";
import type { Ref } from "vue";

type Coorindates = {
  x: Ref<number>;
  y: Ref<number>;
};

export function useTrackCoordinates(el: HTMLElement): Coorindates {
  const x = ref(0);
  const y = ref(0);

  const elOffset = el.getBoundingClientRect();

  x.value = elOffset.left;
  y.value = elOffset.top;

  return { x, y };
}
