#pragma once

#include <memory>
#include <vector>
#include <utility>

namespace ql::parser {
    class AbstractNode : protected std::enable_shared_from_this<AbstractNode> {
    public:
        typedef std::vector<std::shared_ptr<AbstractNode>> ChildrenRef;
        typedef std::weak_ptr<AbstractNode> ParentRef;
    protected:
        ChildrenRef m_Children;
        ParentRef m_Parent;
    public:
        explicit AbstractNode(ParentRef parent) : m_Parent(std::move(parent)) {}

        void addChild(std::shared_ptr<AbstractNode> const& node);
    };
}
