import {
  patchState, signalStore, withMethods, withState,
} from '@ngrx/signals';

type AppState = {
  isOpenActive: boolean;
  isDirectActive: boolean;
  active: string;
};

const initialState: AppState = {
  isOpenActive: false,
  isDirectActive: false,
  active: '',
};

export const ChannelsStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withMethods((store) => ({
    setIsChannelsActive(): void {
      patchState(store, () => ({
        isOpenActive: true,
        isDirectActive: false,
      }));
    },
    setIsContactsActive(): void {
      patchState(store, () => ({
        isOpenActive: true,
        isDirectActive: false,
      }));
    },
    setActiveChannel(channelId: string): void {
      patchState(store, () => ({
        active: channelId,
      }));
    },
  })),
);
