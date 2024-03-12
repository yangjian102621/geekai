import { ref, type Ref } from "vue";
function useState<T>(defaultValue?: T): [Ref<T>, (newValue: T) => void] {
  const state = ref<T>(defaultValue) as Ref<T>;
  const setState = (newValue: T) => {
    state.value = newValue;
  };
  return [state, setState];
}

export default useState;
