import { ButtonHTMLAttributes, FC, ReactNode, Ref } from 'react';

import { VariantProps, tv } from 'shared/lib';

export const buttonVariants = tv({
  base: 'inline-block rounded-md focus-visible:opacity-70 hover:opacity-70 disabled:opacity-60 disabled:cursor-not-allowed',
  variants: {
    variant: {
      clear: 'bg-transparent',
      outlined: 'border-1 border-primary bg-primary',
      accent: 'bg-accent',
    },
    textColor: {
      secondary: 'text-secondary',
    },
    size: {
      md: 'py-[7px] px-[33px]',
      xl: 'py-[13px] px-[33px]',
    },
    icons: {
      true: 'px-[22px] flex items-center justify-center gap-[10px]',
      false: '',
    },
  },
  compoundVariants: [
    { variant: 'clear', class: 'p-0 text-inherit' },
    { variant: 'accent', class: 'text-bg-accent-contrast' },
  ],
  defaultVariants: {
    textColor: 'secondary',
    variant: 'outlined',
    size: 'md',
  },
});

interface ButtonProps
  extends ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  leftIcon?: ReactNode;
  rightIcon?: ReactNode;
  ref?: Ref<HTMLButtonElement>;
}

export const Button: FC<ButtonProps> = ({
  className,
  size,
  variant,
  leftIcon,
  rightIcon,
  ref,
  textColor,
  type = 'button',
  ...props
}) => {
  return (
    <button
      {...props}
      ref={ref}
      className={buttonVariants({
        className,
        variant,
        size,
        textColor,
        icons: !!leftIcon || !!rightIcon,
      })}
      type={type}
    >
      {leftIcon && (
        <span className="flex w-[25px] h-[25px] items-center justify-center">
          {leftIcon}
        </span>
      )}
      {props.children}
      {rightIcon && (
        <span className="flex w-[25px] h-[25px] items-center justify-center">
          {rightIcon}
        </span>
      )}
    </button>
  );
};

Button.displayName = 'Button';
