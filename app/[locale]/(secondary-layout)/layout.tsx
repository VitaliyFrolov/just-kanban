import { FC, PropsWithChildren } from 'react';

const SecondaryLayout: FC<PropsWithChildren> = ({ children }) => {
  return (
    <div className="h-screen flex flex-col bg-[url(/raster-images/registration-bg.webp)] bg-no-repeat bg-cover">
      <main className="shrink grow flex items-center">
        <div className="content-container px-s tablet:px-m">{children}</div>
      </main>
    </div>
  );
};

export default SecondaryLayout;
