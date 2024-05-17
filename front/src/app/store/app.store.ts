import {
  patchState, signalStore, withMethods, withState,
} from '@ngrx/signals';

type AppState = {
  isChannelsActive: boolean;
  isContactsActive: boolean;
};

const initialState: AppState = {
  isChannelsActive: false,
  isContactsActive: false,
};

export const AppStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withMethods((store) => ({
    setIsChannelsActive(active: boolean): void {
      patchState(store, () => ({ isChannelsActive: active, isContactsActive: false }));
    },
    setIsContactsActive(active: boolean): void {
      patchState(store, () => ({ isContactsActive: active, isChannelsActive: false }));
    },
  })),
);
