import React, { useState } from 'react';
import { App, ColorPicker } from 'antd';
const Demo = () => {
  const [value, setValue] = useState('#1677ff');
  return (
    <App>
      <ColorPicker
        value={value}
        onChangeComplete={(colorNew) => {
          setValue(colorNew);
          alert(`The selected color is ${colorNew.toHexString()}`);
          //Post color to /player
          //doPostRequest("/player", colorNew)
        }}
      />
    </App>
  );
};
export default Demo;