import { get } from "@vueuse/core";
import { Item } from "ant-design-vue/es/menu";
import {PermissionsList} from '../config/permissions'


//查找在一个对象数组里面，是否存在一个对象的property=value，有的话返回该对象，没有返回null
export const  FindDuplicates =(items,property,value)=>{

    return items.find(item => item[property] === value);

}

//从localhost里面解析出该用户的权限
export const getUserPermissions = ()=>{

    const result = {};
    
    PermissionsList.forEach(k => result[k] = false);

    const Permissions=localStorage.getItem("UserInfo").permissions

    Permissions.forEach(item => {
        keys.forEach(k => {
          if (item.includes(k)) {
            result[k] = true;
          }
        });
      });

      return result

}
  
