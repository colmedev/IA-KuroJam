
interface LoginButtonProps {
    onClick: () => void;
}

export const Login: React.FC<LoginButtonProps> = ({onClick}) => {

  return (
    <button className="button-login" onClick={onClick}>
        Sign In      
    </button>
  )
}
