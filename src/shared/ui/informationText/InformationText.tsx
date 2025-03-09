import { FC, PropsWithChildren } from 'react';

import { clsx } from 'shared/lib';

type InformationTextProps = PropsWithChildren<{
  className?: string;
  error?: boolean;
}>;

export const InformationText: FC<InformationTextProps> = ({
  className,
  children,
  error,
}) => {
  return (
    <div className={clsx(className, 'text-sm', { 'text-danger': error })}>
      {children}
    </div>
  );
};
