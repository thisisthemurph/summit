import {ReactNode} from "react";

type ContainerProps = {
  children: ReactNode;
}

function Container({ children }: ContainerProps) {
  return (
    <section className="p-4">
      {children}
    </section>
  )
}

export default Container;
