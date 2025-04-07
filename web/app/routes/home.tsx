import X3DHSender from "~/components/x3dh-sender/x3dh-sender.component";
import type { Route } from "./+types/home";
import X3DHReceiver from "~/components/x3dh-receiver/xsdh-receiver.component";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Telegrach" },
    { name: "description", content: "Welcome to Telegrach!" },
  ];
}

export default function Home() {
  return (
    <div>
      <X3DHSender />
      <X3DHReceiver />
    </div>
  );
}
