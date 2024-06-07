import { FaGithub } from "react-icons/fa";
import { FaXTwitter } from "react-icons/fa6";
import { FaLinkedin } from "react-icons/fa6";

type FooterProps = {
  className?: string;
};

export function Footer({ className }: FooterProps) {
  return (
    <footer
      className={`
        ${className}
        flex
        flex-col
        py-6
        gap-6
        border-t-2
      `}
      >
        <address className="not-italic">
          <ul
            className="
              flex
              gap-6
              text-2xl
            ">
            <li className="hover:cursor-pointer">
              <FaGithub title="Visit project on Github." href="https://github.com/buarki/find-castles"/>
            </li>
            <li className="hover:cursor-pointer">
              <FaXTwitter title="Follow on Twitter" href="https://x.com/buarki"/>
            </li>
            <li className="hover:cursor-pointer">
              <FaLinkedin className="text-blue-600" title="Connect on LinkedIn" href="https://www.linkedin.com/in/aurelio-buarque"/>
            </li>
          </ul>
        </address>
      <span>Made with love by <a title="Visit author's website" href="https://buarki.com" className="underline">buarki.com</a></span>
    </footer>
  );
}
