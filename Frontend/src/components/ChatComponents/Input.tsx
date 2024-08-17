import { useState } from "react"
import React from "react"
type InputProps = {
  prompt?: string;
}

export const Input = ({ prompt = '' }: InputProps) =>{

  const [_prompt, setPrompt] = useState(prompt)
  const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPrompt(e.target.value);
    return React.createElement('div', {
      role: "Assistant",
      children: _prompt
    })
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => { 
    e.preventDefault()    
  }
  
  return(
    <form onSubmit={handleSubmit} method="POST">
      <input  placeholder="Type your message here" value={_prompt} onChange={handleInput} className="prompt" />
      <button type="submit" className="button--submit">Send</button>
    </form>
  )

  
}
