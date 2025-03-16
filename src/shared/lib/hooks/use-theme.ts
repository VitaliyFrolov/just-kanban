import { use } from 'react';

import { ThemeContext } from 'shared/context';

export const useTheme = () => use(ThemeContext);
