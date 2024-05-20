import {
  patchState, signalStore, withMethods, withState,
} from '@ngrx/signals';

type AppState = {
  isChannelsActive: boolean;
  isContactsActive: boolean;
  activeChannel: string;
};

const initialState: AppState = {
  isChannelsActive: false,
  isContactsActive: false,
  activeChannel: '',
};

export const AppStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withMethods((store) => ({
    setIsChannelsActive(): void {
      patchState(store, () => ({
        isChannelsActive: true,
        isContactsActive: false,
      }));
    },
    setIsContactsActive(): void {
      patchState(store, () => ({
        isContactsActive: true,
        isChannelsActive: false,
      }));
    },
    setActiveChannel(channelId: string): void {
      patchState(store, () => ({
        activeChannel: channelId,
      }));
    },
  })),
);
