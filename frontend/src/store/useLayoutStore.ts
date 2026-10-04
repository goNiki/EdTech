import { create } from 'zustand';

interface LayoutState {
  isMobileOpen: boolean;
  openMobileMenu: () => void;
  closeMobileMenu: () => void;
  toggleMobileMenu: () => void;
  isCollapsed: boolean;
  toggleCollapsed: () => void;
  setCollapsed: (collapsed: boolean) => void;
}

export const useLayoutStore = create<LayoutState>((set) => ({
  isMobileOpen: false,
  openMobileMenu: () => set({ isMobileOpen: true }),
  closeMobileMenu: () => set({ isMobileOpen: false }),
  toggleMobileMenu: () => set((state) => ({ isMobileOpen: !state.isMobileOpen })),
  isCollapsed: false,
  toggleCollapsed: () => set((state) => ({ isCollapsed: !state.isCollapsed })),
  setCollapsed: (collapsed: boolean) => set({ isCollapsed: collapsed }),
}));
