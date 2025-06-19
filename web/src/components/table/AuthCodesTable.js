import React, { useEffect, useState } from 'react';
import {
  API,
  copy,
  showError,
  showSuccess,
  timestamp2string
} from '../../helpers';

import {
  CheckCircle,
  XCircle,
  Minus,
  HelpCircle,
  User,
  UserCheck,
  Crown,
  Bot,
  Clock
} from 'lucide-react';

import { ITEMS_PER_PAGE } from '../../constants';
import {
  Button,
  Card,
  Divider,
  Dropdown,
  Empty,
  Form,
  Modal,
  Popover,
  Space,
  Table,
  Tag,
  Typography
} from '@douyinfe/semi-ui';
import {
  IllustrationNoResult,
  IllustrationNoResultDark
} from '@douyinfe/semi-illustrations';
import {
  IconPlus,
  IconCopy,
  IconSearch,
  IconEyeOpened,
  IconEdit,
  IconDelete,
  IconStop,
  IconPlay,
  IconMore,
  IconRefresh
} from '@douyinfe/semi-icons';
import EditAuthCode from '../../pages/AuthCode/EditAuthCode';
import BatchCreateAuthCode from '../../pages/AuthCode/BatchCreateAuthCode';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

function renderTimestamp(timestamp) {
  return <>{timestamp2string(timestamp)}</>;
}

const AuthCodesTable = () => {
  const { t } = useTranslation();

  const renderStatus = (status) => {
    switch (status) {
      case 1:
        return (
          <Tag color='green' size='large' shape='circle' prefixIcon={<CheckCircle size={14} />}>
            {t('启用')}
          </Tag>
        );
      case 2:
        return (
          <Tag color='red' size='large' shape='circle' prefixIcon={<XCircle size={14} />}>
            {t('禁用')}
          </Tag>
        );
      case 3:
        return (
          <Tag color='grey' size='large' shape='circle' prefixIcon={<Minus size={14} />}>
            {t('已使用')}
          </Tag>
        );
      case 4:
        return (
          <Tag color='orange' size='large' shape='circle' prefixIcon={<Clock size={14} />}>
            {t('待激活')}
          </Tag>
        );
      case 5:
        return (
          <Tag color='blue' size='large' shape='circle' prefixIcon={<CheckCircle size={14} />}>
            {t('激活')}
          </Tag>
        );
      default:
        return (
          <Tag color='black' size='large' shape='circle' prefixIcon={<HelpCircle size={14} />}>
            {t('未知状态')}
          </Tag>
        );
    }
  };

  const renderUserType = (userType) => {
    switch (userType) {
      case 1:
        return (
          <Tag color='blue' size='large' shape='circle' prefixIcon={<User size={14} />}>
            {t('普通用户')}
          </Tag>
        );
      case 10:
        return (
          <Tag color='orange' size='large' shape='circle' prefixIcon={<UserCheck size={14} />}>
            {t('管理员')}
          </Tag>
        );
      case 100:
        return (
          <Tag color='red' size='large' shape='circle' prefixIcon={<Crown size={14} />}>
            {t('超级管理员')}
          </Tag>
        );
      default:
        return (
          <Tag color='black' size='large' shape='circle' prefixIcon={<HelpCircle size={14} />}>
            {t('未知类型')}
          </Tag>
        );
    }
  };

  const renderExpiredTime = (expiredTime) => {
    if (expiredTime === -1) {
      return (
        <Tag color='green' size='large' shape='circle'>
          {t('永不过期')}
        </Tag>
      );
    }
    
    const now = Math.floor(Date.now() / 1000);
    const isExpired = expiredTime < now;
    
    return (
      <Tag 
        color={isExpired ? 'red' : 'orange'} 
        size='large' 
        shape='circle' 
        prefixIcon={<Clock size={14} />}
      >
        {renderTimestamp(expiredTime)}
      </Tag>
    );
  };

  const columns = [
    {
      title: t('ID'),
      dataIndex: 'id',
      width: 80,
    },
    {
      title: t('授权码'),
      dataIndex: 'code',
      width: 200,
      render: (text, record, index) => {
        return (
          <div className="flex items-center gap-2">
            <Text code copyable={{ content: text }}>
              {text.substring(0, 8)}...
            </Text>
          </div>
        );
      },
    },
    {
      title: t('名称'),
      dataIndex: 'name',
      width: 150,
    },
    {
      title: t('状态'),
      dataIndex: 'status',
      width: 100,
      render: (text, record, index) => {
        return <div>{renderStatus(text)}</div>;
      },
    },
    {
      title: t('用户类型'),
      dataIndex: 'user_type',
      width: 120,
      render: (text, record, index) => {
        return <div>{renderUserType(text)}</div>;
      },
    },
    {
      title: t('过期时间'),
      dataIndex: 'expired_time',
      width: 150,
      render: (text, record, index) => {
        return <div>{renderExpiredTime(text)}</div>;
      },
    },
    {
      title: t('机器人'),
      dataIndex: 'is_bot',
      width: 80,
      render: (text, record, index) => {
        return text ? (
          <Tag color='purple' size='large' shape='circle' prefixIcon={<Bot size={14} />}>
            {t('是')}
          </Tag>
        ) : (
          <Tag color='grey' size='large' shape='circle'>
            {t('否')}
          </Tag>
        );
      },
    },
    {
      title: t('机器码'),
      dataIndex: 'machine_code',
      width: 150,
      render: (text, record, index) => {
        if (!text) {
          return (
            <Tag color='grey' size='large' shape='circle'>
              {t('未绑定')}
            </Tag>
          );
        }
        return (
          <div className="flex items-center gap-2">
            <Text code copyable={{ content: text }}>
              {text.substring(0, 8)}...
            </Text>
          </div>
        );
      },
    },
    {
      title: t('分组'),
      dataIndex: 'group',
      width: 150,
      render: (text, record, index) => {
        if (!text) {
          return (
            <Tag color='grey' size='large' shape='circle'>
              {t('无分组')}
            </Tag>
          );
        }

        // 解析多个分组
        const groups = text.split(',').filter(g => g.trim() !== '');
        if (groups.length === 0) {
          return (
            <Tag color='grey' size='large' shape='circle'>
              {t('无分组')}
            </Tag>
          );
        }

        // 如果只有一个分组，直接显示
        if (groups.length === 1) {
          return (
            <Tag color='blue' size='large' shape='circle'>
              {groups[0].trim()}
            </Tag>
          );
        }

        // 如果有多个分组，显示第一个和数量
        return (
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
            <Tag color='blue' size='large' shape='circle'>
              {groups[0].trim()}
            </Tag>
            {groups.length > 1 && (
              <Tag color='cyan' size='large' shape='circle'>
                +{groups.length - 1}
              </Tag>
            )}
          </div>
        );
      },
    },
    {
      title: t('创建时间'),
      dataIndex: 'created_time',
      width: 150,
      render: (text, record, index) => {
        return <div>{renderTimestamp(text)}</div>;
      },
    },
    {
      title: t('使用者ID'),
      dataIndex: 'used_user_id',
      width: 100,
      render: (text, record, index) => {
        return <div>{text === 0 ? t('无') : text}</div>;
      },
    },
    {
      title: '',
      dataIndex: 'operate',
      fixed: 'right',
      width: 200,
      render: (text, record, index) => {
        // 创建更多操作的下拉菜单项
        const moreMenuItems = [
          {
            node: 'item',
            name: t('删除'),
            icon: <IconDelete />,
            type: 'danger',
            onClick: () => {
              Modal.confirm({
                title: t('确定是否要删除此授权码？'),
                content: t('此修改将不可逆'),
                onOk: () => {
                  manageAuthCode(record.id, 'delete', record).then(() => {
                    removeRecord(record.id);
                  });
                },
              });
            },
          }
        ];

        // 动态添加启用/禁用按钮
        if (record.status === 1) {
          moreMenuItems.push({
            node: 'item',
            name: t('禁用'),
            icon: <IconStop />,
            type: 'warning',
            onClick: () => {
              manageAuthCode(record.id, 'disable', record);
            },
          });
        } else if (record.status === 2) {
          moreMenuItems.push({
            node: 'item',
            name: t('启用'),
            icon: <IconPlay />,
            type: 'secondary',
            onClick: () => {
              manageAuthCode(record.id, 'enable', record);
            },
          });
        }

        return (
          <Space>
            <Popover
              content={
                <div style={{ padding: 10, maxWidth: 300 }}>
                  <div><strong>{t('授权码')}:</strong> {record.code}</div>
                  <div><strong>{t('描述')}:</strong> {record.description || t('无')}</div>
                  <div><strong>{t('WxAutoX码')}:</strong> {record.wx_auto_x_code || t('无')}</div>
                  <div><strong>{t('机器码')}:</strong> {record.machine_code || t('未绑定')}</div>
                  <div><strong>{t('分组')}:</strong> {
                    record.group ?
                      record.group.split(',').filter(g => g.trim() !== '').join(', ') || t('无分组')
                      : t('无分组')
                  }</div>
                </div>
              }
              position='top'
            >
              <Button
                icon={<IconEyeOpened />}
                theme='light'
                type='tertiary'
                size="small"
                className="!rounded-full"
              >
                {t('查看')}
              </Button>
            </Popover>
            <Button
              icon={<IconCopy />}
              theme='light'
              type='secondary'
              size="small"
              className="!rounded-full"
              onClick={async () => {
                await copyText(record.code);
              }}
            >
              {t('复制')}
            </Button>
            <Button
              icon={<IconEdit />}
              theme='light'
              type='tertiary'
              size="small"
              className="!rounded-full"
              onClick={() => {
                console.log('Editing record:', record);
                setEditingAuthCode({ id: record.id });
                setShowEdit(true);
              }}
              disabled={record.status === 3}
            >
              {t('编辑')}
            </Button>
            <Dropdown
              trigger='click'
              position='bottomRight'
              menu={moreMenuItems}
            >
              <Button
                icon={<IconMore />}
                theme='light'
                type='tertiary'
                size="small"
                className="!rounded-full"
              />
            </Dropdown>
          </Space>
        );
      },
    },
  ];

  const [authCodes, setAuthCodes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activePage, setActivePage] = useState(1);
  const [searching, setSearching] = useState(false);
  const [tokenCount, setTokenCount] = useState(ITEMS_PER_PAGE);
  const [selectedKeys, setSelectedKeys] = useState([]);
  const [pageSize, setPageSize] = useState(ITEMS_PER_PAGE);
  const [editingAuthCode, setEditingAuthCode] = useState({
    id: undefined,
  });
  const [showEdit, setShowEdit] = useState(false);
  const [showBatchCreate, setShowBatchCreate] = useState(false);

  // Form 初始值
  const formInitValues = {
    searchKeyword: '',
  };

  // Form API 引用
  const [formApi, setFormApi] = useState(null);

  // 获取表单值的辅助函数
  const getFormValues = () => {
    const formValues = formApi ? formApi.getValues() : {};
    return {
      searchKeyword: formValues.searchKeyword || '',
    };
  };

  const closeEdit = () => {
    setShowEdit(false);
    setTimeout(() => {
      setEditingAuthCode({
        id: undefined,
      });
    }, 500);
  };

  const closeBatchCreate = () => {
    setShowBatchCreate(false);
  };

  const setAuthCodeFormat = (authCodes) => {
    // 确保每个项目都有 key 属性
    const formattedAuthCodes = authCodes.map(item => ({
      ...item,
      key: item.id // 添加 key 属性
    }));
    setAuthCodes(formattedAuthCodes);
  };

  const loadAuthCodes = async (startIdx, pageSize) => {
    const res = await API.get(
      `/api/auth_code/?p=${startIdx}&page_size=${pageSize}`,
    );
    const { success, message, data } = res.data;
    if (success) {
      const newPageData = data.items;
      setActivePage(data.page);
      setTokenCount(data.total);
      setAuthCodeFormat(newPageData);
    } else {
      showError(message);
    }
    setLoading(false);
  };

  const removeRecord = (id) => {
    let newDataSource = [...authCodes];
    let idx = newDataSource.findIndex((data) => data.id === id);

    if (idx > -1) {
      newDataSource.splice(idx, 1);
      setAuthCodes(newDataSource);
    }
  };

  const copyText = async (text) => {
    if (await copy(text)) {
      showSuccess(t('已复制到剪贴板！'));
    } else {
      Modal.error({
        title: t('无法复制到剪贴板，请手动复制'),
        content: text,
        size: 'large'
      });
    }
  };

  const onPaginationChange = (e, { activePage }) => {
    (async () => {
      if (activePage === Math.ceil(authCodes.length / pageSize) + 1) {
        await loadAuthCodes(activePage - 1, pageSize);
      }
      setActivePage(activePage);
    })();
  };

  useEffect(() => {
    loadAuthCodes(0, pageSize)
      .then()
      .catch((reason) => {
        showError(reason);
      });
  }, [pageSize]);

  const refresh = async () => {
    await loadAuthCodes(activePage - 1, pageSize);
  };

  const manageAuthCode = async (id, action, record) => {
    setLoading(true);
    let data = { id };
    let res;
    switch (action) {
      case 'delete':
        res = await API.delete(`/api/auth_code/${id}/`);
        break;
      case 'enable':
        data.status = 1;
        res = await API.put('/api/auth_code/?status_only=true', data);
        break;
      case 'disable':
        data.status = 2;
        res = await API.put('/api/auth_code/?status_only=true', data);
        break;
    }
    const { success, message } = res.data;
    if (success) {
      showSuccess(t('操作成功完成！'));
      let authCode = res.data.data;
      let newAuthCodes = [...authCodes];
      if (action === 'delete') {
        // 删除操作在removeRecord中处理
      } else {
        record.status = authCode.status;
      }
      setAuthCodes(newAuthCodes);
    } else {
      showError(message);
    }
    setLoading(false);
  };

  const searchAuthCodes = async (keyword = null, page, pageSize) => {
    // 如果没有传递keyword参数，从表单获取值
    if (keyword === null) {
      const formValues = getFormValues();
      keyword = formValues.searchKeyword;
    }

    if (keyword === '') {
      await loadAuthCodes(page, pageSize);
      return;
    }
    setSearching(true);
    const res = await API.get(
      `/api/auth_code/search?keyword=${keyword}&p=${page}&page_size=${pageSize}`,
    );
    const { success, message, data } = res.data;
    if (success) {
      const newPageData = data.items;
      setActivePage(data.page);
      setTokenCount(data.total);
      setAuthCodeFormat(newPageData);
    } else {
      showError(message);
    }
    setSearching(false);
  };

  const sortAuthCode = (key) => {
    if (authCodes.length === 0) return;
    setLoading(true);
    let sortedAuthCodes = [...authCodes];
    sortedAuthCodes.sort((a, b) => {
      return ('' + a[key]).localeCompare(b[key]);
    });
    if (sortedAuthCodes[0].id === authCodes[0].id) {
      sortedAuthCodes.reverse();
    }
    setAuthCodes(sortedAuthCodes);
    setLoading(false);
  };

  const handlePageChange = (page) => {
    setActivePage(page);
    const { searchKeyword } = getFormValues();
    if (searchKeyword === '') {
      loadAuthCodes(page, pageSize).then();
    } else {
      searchAuthCodes(searchKeyword, page, pageSize).then();
    }
  };

  // 确保每个数据项都有 key 属性
  let pageData = authCodes.map(item => ({
    ...item,
    key: item.id
  }));


  const rowSelection = {
    onSelect: (record, selected) => { },
    onSelectAll: (selected, selectedRows) => { },
    onChange: (selectedRowKeys, selectedRows) => {
      setSelectedKeys(selectedRows);
    },
  };

  const handleRow = (record, index) => {
    if (record.status !== 1) {
      return {
        style: {
          background: 'var(--semi-color-disabled-border)',
        },
      };
    } else {
      return {};
    }
  };

  const renderHeader = () => (
    <div className="flex flex-col w-full">
      <div className="mb-2">
        <div className="flex items-center text-orange-500">
          <IconEyeOpened className="mr-2" />
          <Text>{t('授权码管理用于控制用户注册和权限分配，支持设置过期时间、用户类型等。')}</Text>
        </div>
      </div>

      <Divider margin="12px" />

      <div className="flex flex-col md:flex-row justify-between items-center gap-4 w-full">
        <div className="flex gap-2 w-full md:w-auto order-2 md:order-1">
          <Button
            theme='light'
            type='primary'
            icon={<IconPlus />}
            className="!rounded-full w-full md:w-auto"
            onClick={() => {
              setEditingAuthCode({
                id: undefined,
              });
              setShowEdit(true);
            }}
          >
            {t('添加授权码')}
          </Button>
          <Button
            theme='light'
            type='secondary'
            icon={<IconPlus />}
            className="!rounded-full w-full md:w-auto"
            onClick={() => {
              setShowBatchCreate(true);
            }}
          >
            {t('批量生成')}
          </Button>
          <Button
            type='warning'
            icon={<IconCopy />}
            className="!rounded-full w-full md:w-auto"
            onClick={async () => {
              if (selectedKeys.length === 0) {
                showError(t('请至少选择一个授权码！'));
                return;
              }
              let codes = '';
              for (let i = 0; i < selectedKeys.length; i++) {
                codes +=
                  selectedKeys[i].name + '    ' + selectedKeys[i].code + '\n';
              }
              await copyText(codes);
            }}
          >
            {t('复制所选授权码到剪贴板')}
          </Button>
          <Button
            theme='light'
            type='tertiary'
            icon={<IconRefresh />}
            className="!rounded-full w-full md:w-auto"
            onClick={refresh}
          >
            {t('刷新')}
          </Button>
        </div>

        <Form
          initValues={formInitValues}
          getFormApi={(api) => setFormApi(api)}
          onSubmit={() => {
            const { searchKeyword } = getFormValues();
            searchAuthCodes(searchKeyword, 1, pageSize).then();
          }}
          className="w-full md:w-auto order-1 md:order-2"
        >
          <div className="flex gap-2 w-full md:w-auto">
            <Form.Input
              field="searchKeyword"
              showClear
              placeholder={t('搜索授权码、名称或描述...')}
              className="w-full md:w-64"
            />
            <Button
              type="primary"
              htmlType="submit"
              icon={<IconSearch />}
              loading={searching}
              className="!rounded-full"
            >
              {t('搜索')}
            </Button>
          </div>
        </Form>
      </div>
    </div>
  );

  return (
    <>
      <Card>
        {renderHeader()}
        <Table
          style={{ marginTop: 5 }}
          columns={columns}
          dataSource={pageData.map((item, index) => ({ ...item, key: item.id || index }))}
          rowKey="key"
          pagination={{
            currentPage: activePage,
            pageSize: pageSize,
            total: tokenCount,
            showSizeChanger: true,
            pageSizeOpts: [10, 20, 50, 100],
            onPageChange: handlePageChange,
            onPageSizeChange: (size) => {
              setPageSize(size);
              setActivePage(1);
            },
          }}
          loading={loading}
          rowSelection={rowSelection}
          onRow={handleRow}
          scroll={{ x: 1200 }}
          empty={
            <Empty
              image={<IllustrationNoResult />}
              darkModeImage={<IllustrationNoResultDark />}
              description={t('暂无数据')}
            />
          }
        />
      </Card>
      <EditAuthCode
        refresh={refresh}
        visible={showEdit}
        editingAuthCode={editingAuthCode}
        onCancel={closeEdit}
      />
      <BatchCreateAuthCode
        refresh={refresh}
        visible={showBatchCreate}
        onCancel={closeBatchCreate}
      />
    </>
  );
};

export default AuthCodesTable;
