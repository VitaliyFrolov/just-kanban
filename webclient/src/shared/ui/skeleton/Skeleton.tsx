import { FC, PropsWithChildren } from 'react';

import { VariantProps, tv } from 'shared/lib';

export const skeletonVariants = tv({
  base: 'bg-skeleton inline-block w-auto h-auto whitespace-nowrap',
  variants: {
    shape: {
      circle: 'rounded-full',
      textLine: 'rounded-m',
      rectangle: 'rounded-s',
    },
    textSize: {
      xs: 'h-[var(--text-xs)]',
      s: 'h-[var(--text-s)]',
      m: 'h-[var(--text-m)]',
      l: 'h-[var(--text-l)]',
      xl: 'h-[var(--text-xl)]',
    },
  },
  defaultVariants: {
    shape: 'textLine',
  },
});

interface SkeletonProps
  extends PropsWithChildren,
    VariantProps<typeof skeletonVariants> {
  className?: string;
  width?: number | string;
  height?: number | string;
}

export const Skeleton: FC<SkeletonProps> = ({
  className,
  height,
  width,
  shape,
  textSize,
  children,
}) => {
  return (
    <span
      className={skeletonVariants({
        className,
        shape,
        textSize,
      })}
      style={{ width, height }}
    >
      <span className="invisible">{children}</span>
    </span>
  );
};
