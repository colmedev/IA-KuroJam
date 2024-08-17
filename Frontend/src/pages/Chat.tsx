import { Input } from "../components";
import './Chat.css'

export const Chat = () => {

  

  const messages = [
    {
      role: "Assistant",
      message: "Hello, what can I do for you?"
    },
    {
      role: "user",
      message: "i need a text about lore ipsum"
    },
    {
      role: "Assistant",
      message: "here it is: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ac sagittis nibh. Curabitur suscipit interdum turpis id fringilla. Pellentesque eu fringilla dui. Ut nunc sapien, lacinia et viverra vitae, feugiat vel metus. Morbi eget lorem eros. Vestibulum sed urna metus. Nam at ligula iaculis, laoreet lectus in, pellentesque neque. Vivamus vel erat non velit convallis eleifend id vitae sapien. Ut venenatis, nisi vitae malesuada porttitor, turpis dui finibus diam, vel finibus turpis ipsum ut libero. Nulla placerat quis libero at consectetur. Integer scelerisque nunc in orci fermentum tempus.'"
    }
  ]


  
  return(
    <div className="Chat">
      <div>
        {
          messages.map((message, index) => (
            <div className="Chat__message" key={index}>
              <p className={message.role === "Assistant" ? "assistant" : "user"}>{message.message}</p>
            </div>
          ))
        }
      </div>
      
      <div>
      <Input />
      </div>
    </div>
  )
}
