import { FC } from 'react';

import { clsx } from 'shared/lib';

import './Loader.css';

interface LoaderProps {
  className?: string;
}

export const Loader: FC<LoaderProps> = ({ className }) => {
  return <span className={clsx(className, 'loader', 'w-[20px] h-[20px]')} />;
};
