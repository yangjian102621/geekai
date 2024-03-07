const get = (origin) => {
  const { properties, required, type } = origin;
  if (type === "object") {
    const array = Object.keys(properties).reduce((prev, name) => {
      return [
        ...prev,
        {
          name,
          ...properties[name],
          required: required.includes(name),
        },
      ];
    }, []);
    return array;
  }
  return [];
}

const set = (tableData) => {
  const properties = tableData.reduce((prev, curr) => {
    if (curr.name) {
      return {
        ...prev,
        [curr.name]: {
          description: curr.description,
          type: curr.type,
        },
      };
    }
    return prev
  }, {});
  const required = tableData.filter((i) => i.required).map((i) => i.name);
  return { properties, required, type: "object" }
}

export default { get, set }