import React from 'react';

import { Panel, PanelGroup, Affix } from 'rsuite';
import 'rsuite/lib/PanelGroup/styles';
import 'rsuite/lib/Panel/styles';

import { NavHeader } from '../NavHeader';
import { RenderData } from '../RenderData';

import './index.css';

const fileFormatMk = "### ELF File Header\n\n|MEMBER|VALUE|\n|---|---|\n|Class|ELFCLASS64|\n|data|ELFDATA2LSB|\n|ByteOrder|LittleEndian|\n|Version|EV_CURRENT|\n|Os/ABI|ELFOSABI_NONE|\n|ABI Version|0|\n|Type|ET_DYN|\n|Machine|EM_X86_64|\n|Version|EV_CURRENT|\n|Entry|0X5b20|\n\n<small>table 1: file \"/bin/ls\" header, for more information, see: \"readelf -h ...\"</small>\n\n### ELF Prog Header\n\n|Type|Offset|FileSize|VirtAddr|MemSize|PhysAddr|Flags|Align|\n|---|---|---|---|---|---|---|---|\n|PT_PHDR|0X40|0X268|0X40|0X268|0X40|PF_R|0X8|\n|PT_INTERP|0X2a8|0X1c|0X2a8|0X1c|0X2a8|PF_R|0X1|\n|PT_LOAD|0X0|0X3510|0X0|0X3510|0X0|PF_R|0X1000|\n|PT_LOAD|0X4000|0X133d1|0X4000|0X133d1|0X4000|PF_X+PF_R|0X1000|\n|PT_LOAD|0X18000|0X8cc0|0X18000|0X8cc0|0X18000|PF_R|0X1000|\n|PT_LOAD|0X20fd0|0X1298|0X21fd0|0X2588|0X21fd0|PF_W+PF_R|0X1000|\n|PT_DYNAMIC|0X21a58|0X200|0X22a58|0X200|0X22a58|PF_W+PF_R|0X8|\n|PT_NOTE|0X2c4|0X44|0X2c4|0X44|0X2c4|PF_R|0X4|\n|PT_GNU_EH_FRAME|0X1d324|0X954|0X1d324|0X954|0X1d324|PF_R|0X4|\n|PT_GNU_STACK|0X0|0X0|0X0|0X0|0X0|PF_W+PF_R|0X10|\n|PT_GNU_RELRO|0X20fd0|0X1030|0X21fd0|0X1030|0X21fd0|PF_R|0X1|\n\n<small>table 2: file \"/bin/ls\" program headers, for more information, see: \"readelf -l ...\"</small>\n\n"

const Container = (props) => {
  return (
    <div className={'container'}>
      <Affix>
        <NavHeader />
      </Affix>
      <div className='content'>
        <PanelGroup>
          <RenderData name={"系统环境"} render_data={"### 系统\n\nlinux\n\n### 平台\n\namd64\n\n### 环境变量\n\n-    =/home/xuanc/文档/CodePratice/ELF/ELFLoaderAnalyser/target/ela-compressed\n- COLORTERM=truecolor\n- DISPLAY=:0\n- LANG=zh_CN.UTF-8\n- LANGUAGE=\n- LC_ADDRESS=zh_CN.UTF-8\n- LC_IDENTIFICATION=zh_CN.UTF-8\n- LC_MEASUREMENT=zh_CN.UTF-8\n- LC_MONETARY=zh_CN.UTF-8\n- LC_NAME=zh_CN.UTF-8\n- LC_NUMERIC=zh_CN.UTF-8\n- LC_PAPER=zh_CN.UTF-8\n- LC_TELEPHONE=zh_CN.UTF-8\n- LC_TIME=zh_CN.UTF-8\n- PATH=/home/xuanc/.gvm/pkgsets/go1.16/global/bin:/home/xuanc/.gvm/gos/go1.16/bin:/home/xuanc/.gvm/pkgsets/go1.16/global/overlay/bin:/home/xuanc/.gvm/bin:/home/xuanc/.gvm/bin:/home/xuanc/.local/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/opt/cxoffice/bin:/usr/lib/jvm/default/bin:/usr/bin/site_perl:/usr/bin/vendor_perl:/usr/bin/core_perl:/home/xuanc/go/bin\n- TERM=xterm-256color\n- XAUTHORITY=/run/user/1000/.mutter-Xwaylandauth.S5CB00\n- LS_COLORS=rs=0:di=01;34:ln=01;36:mh=00:pi=40;33:so=01;35:do=01;35:bd=40;33;01:cd=40;33;01:or=40;31;01:mi=00:su=37;41:sg=30;43:ca=30;41:tw=30;42:ow=34;42:st=37;44:ex=01;32:*.tar=01;31:*.tgz=01;31:*.arc=01;31:*.arj=01;31:*.taz=01;31:*.lha=01;31:*.lz4=01;31:*.lzh=01;31:*.lzma=01;31:*.tlz=01;31:*.txz=01;31:*.tzo=01;31:*.t7z=01;31:*.zip=01;31:*.z=01;31:*.dz=01;31:*.gz=01;31:*.lrz=01;31:*.lz=01;31:*.lzo=01;31:*.xz=01;31:*.zst=01;31:*.tzst=01;31:*.bz2=01;31:*.bz=01;31:*.tbz=01;31:*.tbz2=01;31:*.tz=01;31:*.deb=01;31:*.rpm=01;31:*.jar=01;31:*.war=01;31:*.ear=01;31:*.sar=01;31:*.rar=01;31:*.alz=01;31:*.ace=01;31:*.zoo=01;31:*.cpio=01;31:*.7z=01;31:*.rz=01;31:*.cab=01;31:*.wim=01;31:*.swm=01;31:*.dwm=01;31:*.esd=01;31:*.jpg=01;35:*.jpeg=01;35:*.mjpg=01;35:*.mjpeg=01;35:*.gif=01;35:*.bmp=01;35:*.pbm=01;35:*.pgm=01;35:*.ppm=01;35:*.tga=01;35:*.xbm=01;35:*.xpm=01;35:*.tif=01;35:*.tiff=01;35:*.png=01;35:*.svg=01;35:*.svgz=01;35:*.mng=01;35:*.pcx=01;35:*.mov=01;35:*.mpg=01;35:*.mpeg=01;35:*.m2v=01;35:*.mkv=01;35:*.webm=01;35:*.webp=01;35:*.ogm=01;35:*.mp4=01;35:*.m4v=01;35:*.mp4v=01;35:*.vob=01;35:*.qt=01;35:*.nuv=01;35:*.wmv=01;35:*.asf=01;35:*.rm=01;35:*.rmvb=01;35:*.flc=01;35:*.avi=01;35:*.fli=01;35:*.flv=01;35:*.gl=01;35:*.dl=01;35:*.xcf=01;35:*.xwd=01;35:*.yuv=01;35:*.cgm=01;35:*.emf=01;35:*.ogv=01;35:*.ogx=01;35:*.aac=00;36:*.au=00;36:*.flac=00;36:*.m4a=00;36:*.mid=00;36:*.midi=00;36:*.mka=00;36:*.mp3=00;36:*.mpc=00;36:*.ogg=00;36:*.ra=00;36:*.wav=00;36:*.oga=00;36:*.opus=00;36:*.spx=00;36:*.xspf=00;36:\n- MAIL=/var/mail/root\n- LOGNAME=root\n- USER=root\n- HOME=/root\n- SHELL=/bin/bash\n- SUDO_COMMAND=../target/ela-compressed -e /bin/ls -u xuanc\n- SUDO_USER=xuanc\n- SUDO_UID=1000\n- SUDO_GID=1000\n- QT_AUTO_SCREEN_SCALE_FACTOR=0\n- QT_QPA_PLATFORMTHEME=gnome\n- EDITOR=/usr/bin/nano\n- DOWNGRADE_FROM_ALA=1\n- JAVA_HOME=/usr/lib/jvm/default\n\n\n"} />
          <RenderData name={"文件格式"} render_data={fileFormatMk} />
          <Panel header='加载过程'>
            <p>加载过程</p>
          </Panel>
        </PanelGroup>
      </div>
    </div>
  );
};

export default Container;
